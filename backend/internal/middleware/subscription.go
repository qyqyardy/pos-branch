package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

// SubscriptionOptions controls plan/expiry gating.
// This POS is single-tenant, so we read from store_settings(id=1).
type SubscriptionOptions struct {
	// RequirePlan: "premium" to gate premium-only features.
	// Empty means any plan is allowed.
	RequirePlan string
	// RequireActive blocks access when paid_until < now().
	RequireActive bool
}

func normalizePlan(plan string) string {
	p := strings.ToLower(strings.TrimSpace(plan))
	switch p {
	case "standard", "premium":
		return p
	default:
		return "premium"
	}
}

func isUndefinedColumn(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// 42703 = undefined_column (old schema before subscription columns)
		return string(pqErr.Code) == "42703"
	}
	return false
}

func loadSubscription(db *sql.DB) (plan string, paidUntil time.Time, active bool, err error) {
	// Default: allow everything for dev/old DB until migrations are applied.
	plan = "premium"
	paidUntil = time.Now().AddDate(10, 0, 0)
	active = true

	var (
		p string
		t time.Time
		a bool
	)

	qErr := db.QueryRow(
		`SELECT plan, paid_until, (paid_until >= now()) AS active
		   FROM store_settings
		  WHERE id = 1`,
	).Scan(&p, &t, &a)
	if qErr != nil {
		if qErr == sql.ErrNoRows {
			return plan, paidUntil, active, nil
		}
		if isUndefinedColumn(qErr) {
			return plan, paidUntil, active, nil
		}
		return "", time.Time{}, false, qErr
	}

	plan = normalizePlan(p)
	if t.IsZero() {
		t = paidUntil
	}
	return plan, t, a, nil
}

func RequireSubscription(db *sql.DB, opts SubscriptionOptions) func(http.Handler) http.Handler {
	requiredPlan := normalizePlan(opts.RequirePlan)
	if strings.TrimSpace(opts.RequirePlan) == "" {
		requiredPlan = ""
	}

	bypass := os.Getenv("SUBSCRIPTION_BYPASS") == "true"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bypass {
				next.ServeHTTP(w, r)
				return
			}

			plan, paidUntil, active, err := loadSubscription(db)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}

			if requiredPlan != "" && plan != requiredPlan {
				http.Error(w, "Fitur ini hanya tersedia di plan Premium", 402)
				return
			}

			if opts.RequireActive && !active {
				http.Error(w, "Subscription sudah habis (paid_until: "+paidUntil.Format("2006-01-02")+")", 402)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
