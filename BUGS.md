# Bug List & Improvements

Tracked issues found during code review and API testing. Each item is a self-contained fix.

---

## Bugs

### BUG-1: UpdateExpense returns generic message instead of full expense object
**File:** `backend/internal/service/expense_service.go` (line 44), `backend/internal/handler/expense_handler.go` (line 160-165)
**Severity:** Medium
**Description:**
`UpdateExpense` in the service layer returns `error` only. The handler responds with `{"message": "expense updated"}` instead of the full expense object with joined fields (user_display_name, tag_name, tag_icon). The frontend needs the updated expense object to refresh the UI without a separate fetch.
**Fix:**
- Change `UpdateExpense` return type to `(*models.Expense, error)`
- After `s.expenseRepo.Update(existing)`, re-fetch with `s.expenseRepo.GetByID(expenseID)`
- Update handler to return the expense object in the response

---

### BUG-2: DailyTrend dates may include time component
**File:** `backend/internal/service/analytics_service.go` (line 229)
**Severity:** Medium
**Description:**
The `GetDailyTrend` SQL query uses `SELECT expense_date` directly. Depending on how SQLite stores the date, this can return `"2026-04-14T00:00:00Z"` instead of the clean `"2026-04-14"` format the frontend expects for chart labels.
**Fix:**
Change the SQL to: `SELECT strftime('%Y-%m-%d', expense_date) as date, SUM(amount) as total, COUNT(*) as cnt`

---

### BUG-3: MonthlyComparison change_percent is 0 when last month is 0 but this month > 0
**File:** `backend/internal/service/analytics_service.go` (line 278-279)
**Severity:** Low
**Description:**
When `LastMonth == 0` and `ThisMonth > 0`, the `ChangePercent` stays at `0` because the code only calculates the percentage when `LastMonth > 0`. This is misleading — if you spent nothing last month but spent something this month, it should show 100% increase (or a special indicator).
**Fix:**
Add an `else if mc.ThisMonth > 0` branch that sets `ChangePercent = 100`.

---

### BUG-4: UpdateMember handler returns generic message instead of updated member
**File:** `backend/internal/handler/family_handler.go` (line 174)
**Severity:** Low
**Description:**
`UpdateMember` responds with `{"message": "member updated"}` instead of the updated member object. The frontend has to re-fetch the member list after every update.
**Fix:**
- Change `FamilyService.UpdateMember` to return `(*models.User, error)`
- Re-fetch the user after update and return it in the response

---

## Missing Validation

### VAL-1: No payment_method enum validation on expense endpoints
**File:** `backend/internal/service/expense_service.go` (AddExpense, UpdateExpense)
**Severity:** High
**Description:**
The `payment_method` field accepts any string. The database has a CHECK constraint (`cash`, `online`, `upi`, `card`) that will cause a cryptic SQL error on invalid values instead of a clean validation error.
**Fix:**
Validate that `payment_method` is one of: `cash`, `online`, `upi`, `card` before hitting the database.

---

### VAL-2: No tag_id existence check on expense creation
**File:** `backend/internal/service/expense_service.go` (AddExpense, UpdateExpense)
**Severity:** Medium
**Description:**
If an invalid `tag_id` is passed, the SQL JOIN or foreign key constraint fails with a database error instead of a user-friendly message.
**Fix:**
Check that the tag exists (and belongs to the user's family or is a global tag) before creating/updating the expense.

---

### VAL-3: No expense_date format validation
**File:** `backend/internal/service/expense_service.go` (AddExpense, UpdateExpense)
**Severity:** Medium
**Description:**
The `expense_date` field is only checked for emptiness. An invalid date string like `"not-a-date"` or `"14/04/2026"` will be inserted into SQLite without error but cause incorrect query results for date-range filters and analytics.
**Fix:**
Validate that `expense_date` matches `YYYY-MM-DD` format using `time.Parse("2006-01-02", expenseDate)`.

---

### VAL-4: No amount upper bound or reasonable limit
**File:** `backend/internal/service/expense_service.go` (AddExpense)
**Severity:** Low
**Description:**
Only checks `amount > 0`. A typo like `999999999` or a negative-zero edge case could cause issues. No maximum limit.
**Fix:**
Add a reasonable upper bound check (e.g., amount <= 10,000,000) and ensure amount is not NaN or Inf.

---

### VAL-5: No PIN format validation for family members
**File:** `backend/internal/service/family_service.go` (AddMember, UpdateMember)
**Severity:** High
**Description:**
PINs are meant to be 4-6 digit numeric codes for easy entry by elder family members. Currently any string is accepted as a PIN — including empty-ish whitespace, letters, or very long strings.
**Fix:**
Validate that PIN is 4-6 digits (numeric only) using a regex like `^\d{4,6}$`.

---

### VAL-6: No username format validation
**File:** `backend/internal/service/family_service.go` (CreateFamily, AddMember)
**Severity:** Medium
**Description:**
Usernames have no format restrictions. Spaces, special characters, very long strings, or empty-after-trim strings are all accepted. This can cause login issues and display problems.
**Fix:**
Validate username: alphanumeric + underscores, 3-30 characters, lowercase only. Regex: `^[a-z0-9_]{3,30}$`.

---

### VAL-7: No display_name length validation
**File:** `backend/internal/service/family_service.go` (CreateFamily, AddMember)
**Severity:** Low
**Description:**
Display names have no length limit. A very long name could break UI layouts.
**Fix:**
Validate display_name: 1-50 characters, trimmed of leading/trailing whitespace.

---

### VAL-8: No family name validation
**File:** `backend/internal/service/family_service.go` (CreateFamily)
**Severity:** Low
**Description:**
Family names only check for empty string. No length limit or format validation.
**Fix:**
Validate family name: 1-100 characters, trimmed.

---

### VAL-9: No owner password strength validation
**File:** `backend/internal/service/family_service.go` (CreateFamily)
**Severity:** Medium
**Description:**
Family owner passwords have no minimum length or complexity requirement. A single character password is accepted.
**Fix:**
Require minimum 6 characters for owner passwords.

---

### VAL-10: No page/page_size bounds validation on expense list
**File:** `backend/internal/handler/expense_handler.go` (ListExpenses)
**Severity:** Low
**Description:**
`page_size` has no upper bound. A client could request `page_size=1000000` and get the entire dataset in one query, causing performance issues.
**Fix:**
Cap `page_size` to a maximum of 100. Ensure `page >= 1`.

---

### VAL-11: No custom date range validation on analytics endpoints
**File:** `backend/internal/handler/analytics_handler.go`
**Severity:** Low
**Description:**
When `range=custom` is used, `from` and `to` query params are passed directly to SQL without format validation. Invalid dates could cause unexpected query results.
**Fix:**
Validate that `from` and `to` are valid `YYYY-MM-DD` dates and that `from <= to`.

---

### VAL-12: Tag name validation is minimal
**File:** `backend/internal/handler/tag_handler.go` (CreateTag)
**Severity:** Low
**Description:**
Tag names only check for empty string. No length limit, no duplicate-within-family check at handler level (only relies on DB unique index which gives a generic error).
**Fix:**
Validate tag name: 1-50 characters, trimmed. Return a clear "tag already exists" error message.

---

## Improvements (Non-bug)

### IMP-1: Expense list filter by payment_method should validate the enum
**File:** `backend/internal/handler/expense_handler.go` (ListExpenses)
**Description:**
The `payment_method` filter query param is passed through without validation. An invalid value silently returns no results instead of an error.

### IMP-2: DeleteFamily should cascade-check or warn about existing expenses
**File:** `backend/internal/service/family_service.go` (DeleteFamily)
**Description:**
Deleting a family silently cascades to delete all members and expenses. Should require confirmation or at least return the count of affected records.

### IMP-3: Deactivated members can still use valid JWTs
**File:** `backend/internal/handler/middleware.go`
**Description:**
When a member is deactivated, their existing JWT remains valid until expiry. The auth middleware should check `is_active` status on each request (or at least periodically).
