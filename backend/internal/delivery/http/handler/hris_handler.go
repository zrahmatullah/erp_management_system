package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/delivery/http/middleware"
)

type HRISHandler struct {
	db *pgxpool.Pool
}

func NewHRISHandler(db *pgxpool.Pool) *HRISHandler {
	return &HRISHandler{db: db}
}

func isAuthorizedApprover(user *middleware.Claims) bool {
	if user == nil {
		return false
	}
	role := user.Role
	if role == "Super Admin" || role == "Manager" || role == "HR Admin" || role == "Owner" {
		return true
	}
	for _, p := range user.Permissions {
		if p == "hris:approve" || p == "payroll:approve" {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// 1. EMPLOYEES MANAGEMENT
// -----------------------------------------------------------------------------

// GetEmployees returns list of employees with department, position, and full database fields
func (h *HRISHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT 
			e.id, e.nik, e.first_name, e.last_name, 
			COALESCE(e.email, ''), COALESCE(e.phone, ''),
			COALESCE(d.id::text, ''), COALESCE(d.name, 'Operations') as department,
			COALESCE(p.id::text, ''), COALESCE(p.title, 'Staff') as position,
			COALESCE(e.job_title, p.title, 'Staff') as job_title,
			e.basic_salary, e.employment_type, e.status, 
			e.join_date, e.end_date,
			COALESCE(e.gender, 'Laki-laki') as gender,
			COALESCE(e.marital_status, 'Belum Kawin') as marital_status,
			COALESCE(e.address, '-') as address,
			COALESCE(e.national_id, '-') as national_id,
			COALESCE(e.tax_id, '-') as tax_id,
			COALESCE(e.bank_name, 'BCA') as bank_name,
			COALESCE(e.bank_account, '-') as bank_account,
			COALESCE(e.bank_account_name, concat(e.first_name, ' ', e.last_name)) as bank_account_name,
			COALESCE(e.photo_url, '') as photo_url,
			COALESCE(e.remaining_leave, 12) as remaining_leave,
			COALESCE(e.branch_id::text, '') as branch_id,
			COALESCE(e.user_id::text, '') as user_id
		FROM employees e
		LEFT JOIN departments d ON e.department_id = d.id
		LEFT JOIN positions p ON e.position_id = p.id
		WHERE e.deleted_at IS NULL
		ORDER BY e.nik ASC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, nik, first, last, email, phone, deptID, dept, posID, pos, jobTitle string
		var empType, status, gender, marital, addr, nationalID, taxID, bankName, bankAcc, bankAccName, photo string
		var branchID, userID string
		var salary float64
		var joinDate time.Time
		var endDate *time.Time
		var remainingLeave int

		if err := rows.Scan(
			&id, &nik, &first, &last, &email, &phone,
			&deptID, &dept, &posID, &pos, &jobTitle,
			&salary, &empType, &status, &joinDate, &endDate,
			&gender, &marital, &addr, &nationalID, &taxID,
			&bankName, &bankAcc, &bankAccName, &photo, &remainingLeave,
			&branchID, &userID,
		); err == nil {
			endStr := ""
			if endDate != nil {
				endStr = endDate.Format("2006-01-02")
			}
			fullName := fmt.Sprintf("%s %s", first, last)
			list = append(list, map[string]interface{}{
				"id":                id,
				"nik":               nik,
				"first_name":        first,
				"last_name":         last,
				"name":              fullName,
				"full_name":         fullName,
				"email":             email,
				"phone":             phone,
				"department_id":     deptID,
				"department":        dept,
				"position_id":       posID,
				"position":          pos,
				"position_name":     pos,
				"job_title":         jobTitle,
				"branch_id":         branchID,
				"user_id":           userID,
				"basic_salary":      salary,
				"employment_type":   empType,
				"status":            status,
				"join_date":         joinDate.Format("2006-01-02"),
				"end_date":          endStr,
				"gender":            gender,
				"marital_status":    marital,
				"address":           addr,
				"national_id":       nationalID,
				"tax_id":            taxID,
				"bank_name":         bankName,
				"bank_account":      bankAcc,
				"bank_account_name": bankAccName,
				"photo_url":         photo,
				"remaining_leave":   remainingLeave,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

// CreateEmployee registers a new employee into PostgreSQL
func (h *HRISHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		NIK             string  `json:"nik"`
		FirstName       string  `json:"first_name"`
		LastName        string  `json:"last_name"`
		Email           string  `json:"email"`
		Phone           string  `json:"phone"`
		BranchID        string  `json:"branch_id"`
		DepartmentID    string  `json:"department_id"`
		PositionID      string  `json:"position_id"`
		BasicSalary     float64 `json:"basic_salary"`
		EmploymentType  string  `json:"employment_type"`
		Status          string  `json:"status"`
		Gender          string  `json:"gender"`
		MaritalStatus   string  `json:"marital_status"`
		Address         string  `json:"address"`
		NationalID      string  `json:"national_id"`
		TaxID           string  `json:"tax_id"`
		BankName        string  `json:"bank_name"`
		BankAccount     string  `json:"bank_account"`
		BankAccountName string  `json:"bank_account_name"`
		JoinDate        string  `json:"join_date"`
		RemainingLeave  int     `json:"remaining_leave"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Payload karyawan tidak valid")
		return
	}

	if body.FirstName == "" {
		writeError(w, http.StatusBadRequest, "Nama depan wajib diisi")
		return
	}

	// Auto-generate NIK if empty
	if body.NIK == "" {
		var count int
		_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM employees").Scan(&count)
		body.NIK = fmt.Sprintf("EMP-%03d", count+1)
	}

	// Default branch if empty
	if body.BranchID == "" {
		_ = h.db.QueryRow(ctx, "SELECT id FROM branches LIMIT 1").Scan(&body.BranchID)
	}

	if body.Status == "" {
		body.Status = "active"
	}
	if body.EmploymentType == "" {
		body.EmploymentType = "full_time"
	}
	if body.RemainingLeave <= 0 {
		body.RemainingLeave = 12
	}

	joinDate := time.Now()
	if body.JoinDate != "" {
		if t, err := time.Parse("2006-01-02", body.JoinDate); err == nil {
			joinDate = t
		}
	}

	newID := uuid.New().String()
	query := `
		INSERT INTO employees (
			id, branch_id, department_id, position_id, nik,
			first_name, last_name, email, phone, basic_salary,
			employment_type, status, gender, marital_status, address,
			national_id, tax_id, bank_name, bank_account, bank_account_name,
			join_date, remaining_leave, created_at, updated_at
		) VALUES (
			$1, NULLIF($2, '')::uuid, NULLIF($3, '')::uuid, NULLIF($4, '')::uuid, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20,
			$21, $22, NOW(), NOW()
		)`

	_, err := h.db.Exec(ctx, query,
		newID, body.BranchID, body.DepartmentID, body.PositionID, body.NIK,
		body.FirstName, body.LastName, body.Email, body.Phone, body.BasicSalary,
		body.EmploymentType, body.Status, body.Gender, body.MaritalStatus, body.Address,
		body.NationalID, body.TaxID, body.BankName, body.BankAccount, body.BankAccountName,
		joinDate, body.RemainingLeave,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mendaftarkan karyawan: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"id":      newID,
		"nik":     body.NIK,
		"message": fmt.Sprintf("Karyawan %s %s (%s) berhasil didaftarkan ke sistem", body.FirstName, body.LastName, body.NIK),
	})
}

// UpdateEmployee updates an employee profile in PostgreSQL
func (h *HRISHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	var body struct {
		FirstName       string  `json:"first_name"`
		LastName        string  `json:"last_name"`
		Email           string  `json:"email"`
		Phone           string  `json:"phone"`
		DepartmentID    string  `json:"department_id"`
		PositionID      string  `json:"position_id"`
		BasicSalary     float64 `json:"basic_salary"`
		EmploymentType  string  `json:"employment_type"`
		Status          string  `json:"status"`
		Gender          string  `json:"gender"`
		MaritalStatus   string  `json:"marital_status"`
		Address         string  `json:"address"`
		NationalID      string  `json:"national_id"`
		TaxID           string  `json:"tax_id"`
		BankName        string  `json:"bank_name"`
		BankAccount     string  `json:"bank_account"`
		BankAccountName string  `json:"bank_account_name"`
		RemainingLeave  int     `json:"remaining_leave"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Payload update tidak valid")
		return
	}

	query := `
		UPDATE employees SET
			first_name = COALESCE(NULLIF($1, ''), first_name),
			last_name = COALESCE(NULLIF($2, ''), last_name),
			email = COALESCE(NULLIF($3, ''), email),
			phone = COALESCE(NULLIF($4, ''), phone),
			department_id = COALESCE(NULLIF($5, '')::uuid, department_id),
			position_id = COALESCE(NULLIF($6, '')::uuid, position_id),
			basic_salary = CASE WHEN $7 > 0 THEN $7 ELSE basic_salary END,
			employment_type = COALESCE(NULLIF($8, ''), employment_type),
			status = COALESCE(NULLIF($9, ''), status),
			gender = COALESCE(NULLIF($10, ''), gender),
			marital_status = COALESCE(NULLIF($11, ''), marital_status),
			address = COALESCE(NULLIF($12, ''), address),
			national_id = COALESCE(NULLIF($13, ''), national_id),
			tax_id = COALESCE(NULLIF($14, ''), tax_id),
			bank_name = COALESCE(NULLIF($15, ''), bank_name),
			bank_account = COALESCE(NULLIF($16, ''), bank_account),
			bank_account_name = COALESCE(NULLIF($17, ''), bank_account_name),
			remaining_leave = CASE WHEN $18 >= 0 THEN $18 ELSE remaining_leave END,
			updated_at = NOW()
		WHERE id = $19 AND deleted_at IS NULL`

	_, err := h.db.Exec(ctx, query,
		body.FirstName, body.LastName, body.Email, body.Phone,
		body.DepartmentID, body.PositionID, body.BasicSalary,
		body.EmploymentType, body.Status, body.Gender, body.MaritalStatus,
		body.Address, body.NationalID, body.TaxID, body.BankName,
		body.BankAccount, body.BankAccountName, body.RemainingLeave, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal memperbarui karyawan: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profil karyawan berhasil diperbarui",
	})
}

// DeleteEmployee soft deletes an employee
func (h *HRISHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	_, err := h.db.Exec(ctx, "UPDATE employees SET deleted_at = NOW(), status = 'resigned' WHERE id = $1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Karyawan berhasil dinonaktifkan"})
}

// -----------------------------------------------------------------------------
// 2. ATTENDANCE MANAGEMENT
// -----------------------------------------------------------------------------

// GetAttendances returns logs of clock in/out
func (h *HRISHandler) GetAttendances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dateParam := r.URL.Query().Get("date")

	query := `
		SELECT 
			a.id, e.id as emp_id, e.nik, concat(e.first_name, ' ', e.last_name) as emp_name,
			COALESCE(s.name, 'Pagi') as shift_name, 
			a.clock_in, a.clock_out,
			a.status, a.late_minutes, a.overtime_minutes, COALESCE(a.notes, '')
		FROM attendances a
		JOIN employees e ON a.employee_id = e.id
		LEFT JOIN work_shifts s ON a.shift_id = s.id`

	var args []interface{}
	if dateParam != "" {
		query += " WHERE a.clock_in::date = $1"
		args = append(args, dateParam)
	}
	query += " ORDER BY a.clock_in DESC LIMIT 100"

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, empID, nik, empName, shiftName, status, notes string
		var clockIn time.Time
		var clockOut *time.Time
		var lateMin, otMin int

		if err := rows.Scan(&id, &empID, &nik, &empName, &shiftName, &clockIn, &clockOut, &status, &lateMin, &otMin, &notes); err == nil {
			clockOutStr := "-"
			if clockOut != nil {
				clockOutStr = clockOut.Format("15:04")
			}
			list = append(list, map[string]interface{}{
				"id":               id,
				"employee_id":      empID,
				"nik":              nik,
				"name":             empName,
				"shift":            shiftName,
				"clock_in":         clockIn.Format("15:04"),
				"clock_out":        clockOutStr,
				"date":             clockIn.Format("2006-01-02"),
				"status":           status,
				"late_minutes":     lateMin,
				"overtime_minutes": otMin,
				"notes":            notes,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

// GetTodayAttendanceStatus returns current attendance summary and employee's clock status
func (h *HRISHandler) GetTodayAttendanceStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	today := time.Now().Format("2006-01-02")

	// Get counts for today
	var totalEmployees, presentCount, lateCount, leaveCount int
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM employees WHERE status = 'active' AND deleted_at IS NULL").Scan(&totalEmployees)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM attendances WHERE clock_in::date = $1 AND status = 'present'", today).Scan(&presentCount)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM attendances WHERE clock_in::date = $1 AND status = 'late'", today).Scan(&lateCount)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM leaves WHERE status = 'approved' AND $1::date BETWEEN start_date AND end_date", today).Scan(&leaveCount)

	// Check if current user or specified employee has clocked in today
	empID := r.URL.Query().Get("employee_id")
	hasClockedIn := false
	var clockInTime, clockOutTime, attendanceID string

	if empID != "" {
		var inTime time.Time
		var outTime *time.Time
		err := h.db.QueryRow(ctx, `
			SELECT id, clock_in, clock_out 
			FROM attendances 
			WHERE employee_id = $1 AND clock_in::date = $2 
			ORDER BY clock_in DESC LIMIT 1`, empID, today).Scan(&attendanceID, &inTime, &outTime)
		if err == nil {
			hasClockedIn = true
			clockInTime = inTime.Format("15:04")
			if outTime != nil {
				clockOutTime = outTime.Format("15:04")
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"date":            today,
		"total_employees": totalEmployees,
		"total_present":   presentCount + lateCount,
		"on_time":         presentCount,
		"late":            lateCount,
		"on_leave":        leaveCount,
		"has_clocked_in":  hasClockedIn,
		"attendance_id":   attendanceID,
		"clock_in_time":   clockInTime,
		"clock_out_time":  clockOutTime,
	})
}

// ClockIn records an employee clock-in
func (h *HRISHandler) ClockIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		EmployeeID string `json:"employee_id"`
		Notes      string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.EmployeeID == "" {
		// Fallback to first active employee
		_ = h.db.QueryRow(ctx, "SELECT id FROM employees WHERE status = 'active' AND deleted_at IS NULL LIMIT 1").Scan(&body.EmployeeID)
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	// Check if already clocked in today
	var existingID string
	err := h.db.QueryRow(ctx, "SELECT id FROM attendances WHERE employee_id = $1 AND clock_in::date = $2", body.EmployeeID, today).Scan(&existingID)
	if err == nil {
		writeError(w, http.StatusBadRequest, "Karyawan ini telah melakukan Clock-In untuk hari ini")
		return
	}

	// Check schedule shift for today
	var shiftID *string
	var shiftStartTime *string
	_ = h.db.QueryRow(ctx, `
		SELECT es.shift_id::text, ws.start_time::text
		FROM employee_schedules es
		JOIN work_shifts ws ON es.shift_id = ws.id
		WHERE es.employee_id = $1 AND es.date = $2`, body.EmployeeID, today).Scan(&shiftID, &shiftStartTime)

	status := "present"
	lateMinutes := 0

	// Determine if late based on shift start time (default 07:00)
	targetHour := 7
	targetMin := 0
	if shiftStartTime != nil && *shiftStartTime != "" {
		_, _ = fmt.Sscanf(*shiftStartTime, "%02d:%02d", &targetHour, &targetMin)
	}

	shiftTimeToday := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMin, 0, 0, now.Location())
	if now.After(shiftTimeToday.Add(5 * time.Minute)) {
		status = "late"
		lateMinutes = int(now.Sub(shiftTimeToday).Minutes())
	}

	attID := uuid.New().String()
	_, err = h.db.Exec(ctx, `
		INSERT INTO attendances (id, employee_id, shift_id, clock_in, status, late_minutes, notes, created_at)
		VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, NOW())`,
		attID, body.EmployeeID, shiftID, now, status, lateMinutes, body.Notes,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mencatat Clock In: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":      true,
		"id":           attID,
		"status":       status,
		"late_minutes": lateMinutes,
		"clock_in":     now.Format("15:04:05"),
		"message":      fmt.Sprintf("Clock In berhasil dicatat pada %s (%s)", now.Format("15:04"), status),
	})
}

// ClockOut records clock out time and calculates hours worked
func (h *HRISHandler) ClockOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		EmployeeID string `json:"employee_id"`
		Notes      string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.EmployeeID == "" {
		_ = h.db.QueryRow(ctx, "SELECT id FROM employees WHERE status = 'active' AND deleted_at IS NULL LIMIT 1").Scan(&body.EmployeeID)
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	var attID string
	var clockIn time.Time
	err := h.db.QueryRow(ctx, `
		SELECT id, clock_in 
		FROM attendances 
		WHERE employee_id = $1 AND clock_in::date = $2 
		ORDER BY clock_in DESC LIMIT 1`, body.EmployeeID, today).Scan(&attID, &clockIn)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Belum ada catatan Clock In untuk hari ini. Silakan Clock In terlebih dahulu.")
		return
	}

	// Calculate overtime if worked more than 8 hours
	duration := now.Sub(clockIn)
	otMinutes := 0
	if duration.Hours() > 8.0 {
		otMinutes = int((duration - (8 * time.Hour)).Minutes())
	}

	notesUpdate := body.Notes
	_, err = h.db.Exec(ctx, `
		UPDATE attendances 
		SET clock_out = $1, overtime_minutes = $2, notes = COALESCE(NULLIF($3, ''), notes)
		WHERE id = $4`, now, otMinutes, notesUpdate, attID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mencatat Clock Out: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":          true,
		"clock_out":        now.Format("15:04:05"),
		"overtime_minutes": otMinutes,
		"duration_hours":   math.Round(duration.Hours()*10) / 10,
		"message":          fmt.Sprintf("Clock Out berhasil dicatat pada %s. Terima kasih atas kerja keras Anda!", now.Format("15:04")),
	})
}

// -----------------------------------------------------------------------------
// 3. LEAVE MANAGEMENT & APPROVAL FLOW
// -----------------------------------------------------------------------------

// GetLeaves returns list of leave applications
func (h *HRISHandler) GetLeaves(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT 
			l.id, e.id as emp_id, e.nik, concat(e.first_name, ' ', e.last_name) as emp_name,
			l.leave_type, l.start_date, l.end_date, l.total_days, l.reason, l.status,
			COALESCE(u.username, '') as approved_by_name, l.created_at
		FROM leaves l
		JOIN employees e ON l.employee_id = e.id
		LEFT JOIN users u ON l.approved_by = u.id
		WHERE l.deleted_at IS NULL
		ORDER BY l.created_at DESC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, empID, nik, name, lType, reason, status, approver string
		var start, end, createdAt time.Time
		var days int

		if err := rows.Scan(&id, &empID, &nik, &name, &lType, &start, &end, &days, &reason, &status, &approver, &createdAt); err == nil {
			list = append(list, map[string]interface{}{
				"id":               id,
				"employee_id":      empID,
				"nik":              nik,
				"name":             name,
				"leave_type":       lType,
				"start_date":       start.Format("2006-01-02"),
				"end_date":         end.Format("2006-01-02"),
				"total_days":       days,
				"reason":           reason,
				"status":           status,
				"approved_by_name": approver,
				"created_at":       createdAt.Format("2006-01-02 15:04"),
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

// CreateLeave submits a new leave request
func (h *HRISHandler) CreateLeave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		EmployeeID string `json:"employee_id"`
		LeaveType  string `json:"leave_type"` // 'Cuti Tahunan', 'Cuti Sakit', 'Cuti Alasan Penting'
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		TotalDays  int    `json:"total_days"`
		Reason     string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid leave payload")
		return
	}

	if body.EmployeeID == "" || body.StartDate == "" || body.EndDate == "" {
		writeError(w, http.StatusBadRequest, "Employee ID, tanggal mulai, dan tanggal selesai wajib diisi")
		return
	}

	start, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Format start_date harus YYYY-MM-DD")
		return
	}
	end, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Format end_date harus YYYY-MM-DD")
		return
	}

	if end.Before(start) {
		writeError(w, http.StatusBadRequest, "Tanggal selesai tidak boleh sebelum tanggal mulai")
		return
	}

	if body.TotalDays <= 0 {
		body.TotalDays = int(end.Sub(start).Hours()/24) + 1
	}

	if body.LeaveType == "" {
		body.LeaveType = "Cuti Tahunan"
	}

	newID := uuid.New().String()
	query := `
		INSERT INTO leaves (id, employee_id, leave_type, start_date, end_date, total_days, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', NOW(), NOW())`

	_, err = h.db.Exec(ctx, query, newID, body.EmployeeID, body.LeaveType, start, end, body.TotalDays, body.Reason)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengajukan cuti: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"id":      newID,
		"status":  "pending",
		"message": fmt.Sprintf("Pengajuan %s (%d hari) berhasil diajukan dan menunggu persetujuan Manager.", body.LeaveType, body.TotalDays),
	})
}

// UpdateLeaveStatus handles APPROVAL FLOW by Managers/Super Admins
func (h *HRISHandler) UpdateLeaveStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	user, _ := middleware.GetUserFromContext(ctx)
	if !isAuthorizedApprover(user) {
		writeError(w, http.StatusForbidden, "Hanya Manager, HR Admin, atau Super Admin yang berhak menyetujui/menolak cuti.")
		return
	}

	var body struct {
		Status string `json:"status"` // 'approved' or 'rejected'
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if body.Status != "approved" && body.Status != "rejected" {
		writeError(w, http.StatusBadRequest, "Status harus 'approved' atau 'rejected'")
		return
	}

	// Fetch leave details
	var empID, leaveType string
	var totalDays int
	err := h.db.QueryRow(ctx, "SELECT employee_id, leave_type, total_days FROM leaves WHERE id = $1 AND deleted_at IS NULL", id).Scan(&empID, &leaveType, &totalDays)
	if err != nil {
		writeError(w, http.StatusNotFound, "Permohonan cuti tidak ditemukan")
		return
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	var approverID *uuid.UUID
	if user != nil {
		approverID = &user.UserID
	}

	_, err = tx.Exec(ctx, `
		UPDATE leaves 
		SET status = $1, approved_by = $2, updated_at = NOW()
		WHERE id = $3`, body.Status, approverID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Deduct leave balance if annual leave and approved
	if body.Status == "approved" && leaveType == "Cuti Tahunan" {
		_, _ = tx.Exec(ctx, "UPDATE employees SET remaining_leave = GREATEST(0, remaining_leave - $1) WHERE id = $2", totalDays, empID)
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msg := fmt.Sprintf("Permohonan cuti berhasil disetujui oleh %s", user.Email)
	if body.Status == "rejected" {
		msg = fmt.Sprintf("Permohonan cuti ditolak oleh %s", user.Email)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  body.Status,
		"message": msg,
	})
}

// -----------------------------------------------------------------------------
// 4. SHIFT SCHEDULES (ROSTER PERSISTENCE)
// -----------------------------------------------------------------------------

// GetShiftSchedules returns the persistent weekly shift schedule
func (h *HRISHandler) GetShiftSchedules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDateParam := r.URL.Query().Get("start_date")
	endDateParam := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	if startDateParam != "" {
		startDate, _ = time.Parse("2006-01-02", startDateParam)
	}
	if endDateParam != "" {
		endDate, _ = time.Parse("2006-01-02", endDateParam)
	}

	// Default to current week Monday to Sunday
	if startDate.IsZero() {
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday as 7th day
		}
		startDate = now.AddDate(0, 0, -(weekday - 1))
	}
	if endDate.IsZero() {
		endDate = startDate.AddDate(0, 0, 6)
	}

	// Fetch all active employees
	empRows, err := h.db.Query(ctx, `
		SELECT e.id, e.nik, concat(e.first_name, ' ', e.last_name) as emp_name,
		       COALESCE(d.name, 'Operations') as dept, COALESCE(p.title, 'Staff') as pos
		FROM employees e
		LEFT JOIN departments d ON e.department_id = d.id
		LEFT JOIN positions p ON e.position_id = p.id
		WHERE e.deleted_at IS NULL
		ORDER BY e.nik ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer empRows.Close()

	type empInfo struct {
		id   string
		nik  string
		name string
		dept string
		pos  string
	}
	var emps []empInfo
	for empRows.Next() {
		var e empInfo
		if err := empRows.Scan(&e.id, &e.nik, &e.name, &e.dept, &e.pos); err == nil {
			emps = append(emps, e)
		}
	}

	// Fetch schedules for the date range
	schedQuery := `
		SELECT es.employee_id, es.date, COALESCE(ws.name, 'OFF') as shift_name
		FROM employee_schedules es
		LEFT JOIN work_shifts ws ON es.shift_id = ws.id
		WHERE es.date BETWEEN $1 AND $2`

	schedRows, err := h.db.Query(ctx, schedQuery, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	schedMap := make(map[string]string) // key: employee_id + "_" + dateStr -> shiftName
	if err == nil {
		defer schedRows.Close()
		for schedRows.Next() {
			var eID, sName string
			var d time.Time
			if err := schedRows.Scan(&eID, &d, &sName); err == nil {
				key := fmt.Sprintf("%s_%s", eID, d.Format("2006-01-02"))
				schedMap[key] = sName
			}
		}
	}

	// Build 7-day schedule array for each employee
	var roster []map[string]interface{}
	for _, e := range emps {
		shifts := make([]string, 7)
		for dayIdx := 0; dayIdx < 7; dayIdx++ {
			dayDate := startDate.AddDate(0, 0, dayIdx).Format("2006-01-02")
			key := fmt.Sprintf("%s_%s", e.id, dayDate)
			shiftVal := "OFF"
			if s, exists := schedMap[key]; exists {
				shiftVal = s
			} else {
				// Default alternating shifts if not set yet
				if dayIdx == 6 {
					shiftVal = "OFF"
				} else if dayIdx%2 == 0 {
					shiftVal = "Pagi"
				} else {
					shiftVal = "Siang"
				}
			}
			shifts[dayIdx] = shiftVal
		}

		roster = append(roster, map[string]interface{}{
			"id":         e.id,
			"nik":        e.nik,
			"name":       e.name,
			"department": e.dept,
			"position":   e.pos,
			"shifts":     shifts,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"data":       roster,
	})
}

// SetShiftSchedule assigns a shift to an employee for a specific date (Upsert)
func (h *HRISHandler) SetShiftSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		EmployeeID string `json:"employee_id"`
		ShiftID    string `json:"shift_id"` // can be empty or null for 'OFF'
		ShiftName  string `json:"shift_name"`
		Date       string `json:"date"`
		Notes      string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid schedule payload")
		return
	}

	if body.EmployeeID == "" || body.Date == "" {
		writeError(w, http.StatusBadRequest, "employee_id dan date wajib diisi")
		return
	}

	// Resolve shift_id from shift_name if not provided
	if body.ShiftID == "" && body.ShiftName != "" && body.ShiftName != "OFF" {
		_ = h.db.QueryRow(ctx, "SELECT id::text FROM work_shifts WHERE name = $1 LIMIT 1", body.ShiftName).Scan(&body.ShiftID)
	}

	var shiftUUID *uuid.UUID
	if body.ShiftID != "" && body.ShiftID != "OFF" {
		if u, err := uuid.Parse(body.ShiftID); err == nil {
			shiftUUID = &u
		}
	}

	query := `
		INSERT INTO employee_schedules (employee_id, shift_id, date, notes, created_at, updated_at)
		VALUES ($1::uuid, $2, $3::date, $4, NOW(), NOW())
		ON CONFLICT (employee_id, date) DO UPDATE 
		SET shift_id = EXCLUDED.shift_id, notes = EXCLUDED.notes, updated_at = NOW()`

	_, err := h.db.Exec(ctx, query, body.EmployeeID, shiftUUID, body.Date, body.Notes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal menyimpan jadwal: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Jadwal shift berhasil disimpan untuk tanggal %s", body.Date),
	})
}

// -----------------------------------------------------------------------------
// 5. PAYROLL RUN & ACCOUNTING INTEGRATION
// -----------------------------------------------------------------------------

// GetPayrolls returns payroll records
func (h *HRISHandler) GetPayrolls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	period := r.URL.Query().Get("period")

	query := `
		SELECT 
			pr.id, e.id as emp_id, e.nik, concat(e.first_name, ' ', e.last_name) as emp_name,
			COALESCE(pos.title, 'Barista') as position, COALESCE(dept.name, 'FOH') as department,
			COALESCE(e.bank_name, 'BCA') as bank_name, COALESCE(e.bank_account, '-') as bank_account,
			pr.basic_salary, pr.allowances, pr.overtime_pay, pr.gross_salary,
			pr.bpjs_deduction, pr.tax_deduction, pr.other_deductions, pr.total_deductions, pr.net_salary,
			pr.is_paid, pr.period_start, pr.period_end, pr.paid_at
		FROM payroll_records pr
		JOIN employees e ON pr.employee_id = e.id
		LEFT JOIN positions pos ON e.position_id = pos.id
		LEFT JOIN departments dept ON e.department_id = dept.id
		WHERE pr.deleted_at IS NULL`

	var args []interface{}
	if period != "" {
		query += " AND pr.period_start::text LIKE $1"
		args = append(args, period+"%")
	}
	query += " ORDER BY e.nik ASC"

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, empID, nik, name, pos, dept, bankName, bankAcc string
		var basic, allow, ot, gross, bpjs, tax, otherDed, totalDed, net float64
		var isPaid bool
		var pStart, pEnd time.Time
		var paidAt *time.Time

		if err := rows.Scan(
			&id, &empID, &nik, &name, &pos, &dept, &bankName, &bankAcc,
			&basic, &allow, &ot, &gross, &bpjs, &tax, &otherDed, &totalDed, &net,
			&isPaid, &pStart, &pEnd, &paidAt,
		); err == nil {
			paidStr := "-"
			if paidAt != nil {
				paidStr = paidAt.Format("2006-01-02 15:04")
			}
			list = append(list, map[string]interface{}{
				"id":               id,
				"employee_id":      empID,
				"nik":              nik,
				"name":             name,
				"position":         pos,
				"department":       dept,
				"bank":             fmt.Sprintf("%s - %s", bankName, bankAcc),
				"basicSalary":      basic,
				"allowance":        allow,
				"overtime":         ot,
				"gross":            gross,
				"bpjs":             bpjs,
				"pph21":            tax,
				"other_deductions": otherDed,
				"deduction":        totalDed,
				"net":              net,
				"is_paid":          isPaid,
				"period":           fmt.Sprintf("%s s.d %s", pStart.Format("02 Jan 2006"), pEnd.Format("02 Jan 2006")),
				"period_start":     pStart.Format("2006-01-02"),
				"period_end":       pEnd.Format("2006-01-02"),
				"paid_at":          paidStr,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

// RunPayroll generates or updates payroll calculation for all active employees for the selected period
func (h *HRISHandler) RunPayroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Month int `json:"month"` // 1-12
		Year  int `json:"year"`  // e.g. 2026
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if body.Year == 0 {
		body.Year = time.Now().Year()
	}
	if body.Month == 0 {
		body.Month = int(time.Now().Month())
	}

	firstDay := time.Date(body.Year, time.Month(body.Month), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Fetch all active employees
	empRows, err := h.db.Query(ctx, `
		SELECT id, nik, first_name, last_name, basic_salary
		FROM employees
		WHERE status IN ('active', 'probation') AND deleted_at IS NULL
		ORDER BY nik ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer empRows.Close()

	type empSal struct {
		id     string
		nik    string
		name   string
		salary float64
	}
	var emps []empSal
	for empRows.Next() {
		var e empSal
		var first, last string
		if err := empRows.Scan(&e.id, &e.nik, &first, &last, &e.salary); err == nil {
			e.name = first + " " + last
			emps = append(emps, e)
		}
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	processedCount := 0
	var totalGross, totalNet float64

	for _, e := range emps {
		// Standard cafe allowances (Uang Makan & Transportasi ~ 15-20% basic salary)
		allowance := math.Round(e.salary * 0.18)

		// Overtime from attendances this month
		var otMinutes int
		_ = tx.QueryRow(ctx, `
			SELECT COALESCE(SUM(overtime_minutes), 0)
			FROM attendances
			WHERE employee_id = $1 AND clock_in::date BETWEEN $2 AND $3`,
			e.id, firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")).Scan(&otMinutes)

		otHourlyRate := (e.salary / 173.0) * 1.5
		overtimePay := math.Round((float64(otMinutes) / 60.0) * otHourlyRate)

		grossSalary := e.salary + allowance + overtimePay

		// BPJS Ketenagakerjaan & Kesehatan (approx 3% & 1% employee contribution)
		bpjs := math.Round(e.salary * 0.04)

		// PPh 21 (Simplified PTKP progressive tier for SME Cafe)
		tax := 0.0
		if grossSalary > 4500000 {
			taxable := grossSalary - 4500000
			tax = math.Round(taxable * 0.05)
		}

		totalDeductions := bpjs + tax
		netSalary := grossSalary - totalDeductions

		totalGross += grossSalary
		totalNet += netSalary

		// Upsert into payroll_records
		query := `
			INSERT INTO payroll_records (
				id, employee_id, period_start, period_end, basic_salary,
				allowances, overtime_pay, gross_salary, bpjs_deduction,
				tax_deduction, other_deductions, total_deductions, net_salary,
				is_paid, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9,
				$10, 0, $11, $12,
				FALSE, NOW(), NOW()
			)
			ON CONFLICT (id) DO NOTHING`

		recID := uuid.New().String()
		// Check if record exists for this employee and period
		var existingRecID string
		chkErr := tx.QueryRow(ctx, `
			SELECT id FROM payroll_records 
			WHERE employee_id = $1 AND period_start = $2 AND period_end = $3 AND deleted_at IS NULL`,
			e.id, firstDay, lastDay).Scan(&existingRecID)

		if chkErr == nil {
			// Update existing record
			_, _ = tx.Exec(ctx, `
				UPDATE payroll_records SET
					basic_salary = $1, allowances = $2, overtime_pay = $3, gross_salary = $4,
					bpjs_deduction = $5, tax_deduction = $6, total_deductions = $7, net_salary = $8,
					updated_at = NOW()
				WHERE id = $9 AND is_paid = FALSE`,
				e.salary, allowance, overtimePay, grossSalary, bpjs, tax, totalDeductions, netSalary, existingRecID)
		} else {
			_, _ = tx.Exec(ctx, query,
				recID, e.id, firstDay, lastDay, e.salary,
				allowance, overtimePay, grossSalary, bpjs,
				tax, totalDeductions, netSalary)
		}
		processedCount++
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":         true,
		"period":          fmt.Sprintf("%s - %s", firstDay.Format("02 Jan 2006"), lastDay.Format("02 Jan 2006")),
		"employees_count": processedCount,
		"total_gross":     totalGross,
		"total_net":       totalNet,
		"message":         fmt.Sprintf("Kalkulasi payroll berhasil diproses untuk %d staf. Total Take-Home-Pay: Rp %s", processedCount, formatRupiah(totalNet)),
	})
}

// BatchApprovePayroll approves and marks all payroll records for period as paid, creating journal entry
func (h *HRISHandler) BatchApprovePayroll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, _ := middleware.GetUserFromContext(ctx)
	if !isAuthorizedApprover(user) {
		writeError(w, http.StatusForbidden, "Hanya Super Admin, Manager, atau HR Admin yang berhak menyetujui pencairan payroll.")
		return
	}

	var body struct {
		PeriodStart string `json:"period_start"`
		PeriodEnd   string `json:"period_end"`
		PaymentDate string `json:"payment_date"`
		BankName    string `json:"bank_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	// Fetch total net salary to pay
	var totalNet float64
	var count int
	query := `
		SELECT COUNT(*), COALESCE(SUM(net_salary), 0)
		FROM payroll_records
		WHERE is_paid = FALSE AND deleted_at IS NULL`
	var args []interface{}
	if body.PeriodStart != "" {
		query += " AND period_start >= $1"
		args = append(args, body.PeriodStart)
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&count, &totalNet)
	if err != nil || count == 0 {
		writeError(w, http.StatusBadRequest, "Tidak ada data payroll pending yang siap disetujui")
		return
	}

	// Update payroll records to is_paid = TRUE
	updateQuery := `
		UPDATE payroll_records 
		SET is_paid = TRUE, paid_at = NOW(), updated_at = NOW() 
		WHERE is_paid = FALSE AND deleted_at IS NULL`
	if body.PeriodStart != "" {
		updateQuery += " AND period_start >= $1"
	}
	_, err = tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Post Double-Entry Balanced Accounting Journal Entry
	// Debit: Beban Gaji & Upah (Account 6-1001 or 6100)
	// Credit: Bank BCA Operasional (Account 1-1102)
	var expenseAccountID, bankAccountID string
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE account_type = 'expense' LIMIT 1").Scan(&expenseAccountID)
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '1102' OR account_type = 'asset' LIMIT 1").Scan(&bankAccountID)

	var branchID string
	_ = tx.QueryRow(ctx, "SELECT id FROM branches LIMIT 1").Scan(&branchID)

	if expenseAccountID != "" && bankAccountID != "" {
		journalID := uuid.New().String()
		journalRef := fmt.Sprintf("JRN-PAYROLL-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)

		var createdBy *uuid.UUID
		if user != nil {
			createdBy = &user.UserID
		}

		journalQuery := `
			INSERT INTO journal_entries (id, branch_id, reference_number, entry_date, description, status, created_by, created_at, updated_at)
			VALUES ($1, $2::uuid, $3, NOW()::date, $4, 'posted', $5, NOW(), NOW())`
		desc := fmt.Sprintf("Pencairan Gaji Staf Periode %s (Total %d Karyawan)", body.PeriodStart, count)
		_, jErr := tx.Exec(ctx, journalQuery, journalID, branchID, journalRef, desc, createdBy)

		if jErr == nil {
			// Line 1: Debit Expense
			line1ID := uuid.New().String()
			_, _ = tx.Exec(ctx, `
				INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
				VALUES ($1, $2, $3, 'Beban Gaji Karyawan', $4, 0, NOW())`,
				line1ID, journalID, expenseAccountID, totalNet)

			// Line 2: Credit Bank
			line2ID := uuid.New().String()
			_, _ = tx.Exec(ctx, `
				INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
				VALUES ($1, $2, $3, 'Kas/Bank Pencairan Gaji', 0, $4, NOW())`,
				line2ID, journalID, bankAccountID, totalNet)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":         true,
		"approved_count":  count,
		"total_disbursed": totalNet,
		"approved_by":     user.Email,
		"message":         fmt.Sprintf("Payroll batch (%d staf) berhasil disetujui & dicairkan! Total: Rp %s. Jurnal pengeluaran otomatis dibukukan.", count, formatRupiah(totalNet)),
	})
}

// formatRupiah helper
func formatRupiah(val float64) string {
	return fmt.Sprintf("%.0f", val)
}
