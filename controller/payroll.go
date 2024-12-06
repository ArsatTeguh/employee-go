package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PayrollController struct {
	DB *gorm.DB
}

func (p *PayrollController) GetAll(ctx *gin.Context) {
	err := helper.Premission(ctx)
	if err != nil {
		return
	}

	var payroll []models.Payroll

	search := ctx.Query("date")
	employee_id := ctx.Query("employee_id")
	id, _ := strconv.ParseInt(employee_id, 10, 64)

	var totalCount int64

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	sizePage, _ := strconv.Atoi((ctx.Copy().DefaultQuery("sizePage", "10")))
	offset := (page - 1) * sizePage

	query := p.DB.Model(&payroll).Preload("Employee")

	if search != "" {
		if matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, search); !matched {
			response := &helper.WithoutData{
				Code:    400,
				Message: "invalid date format. Use YYYY-MM",
			}
			response.Response(ctx)
			return
		}
		query = query.Where("DATE_FORMAT(created_at, '%Y-%m') = ?", search)
	}

	if employee_id != "" {
		query = query.Where("employee_id = ?", id)
	}

	query.Count(&totalCount)
	query.Offset(offset).Limit(sizePage).Find(&payroll)

	if len(payroll) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data empty",
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Success Get Payroll",
		Data: map[string]any{
			"payroll":    payroll,                                              // data
			"totalAll":   totalCount,                                           // total data all page
			"total":      len(payroll),                                         // total data per page
			"page":       page,                                                 // current page
			"sizePage":   sizePage,                                             // maximum data per page
			"totalPages": (totalCount + int64(sizePage) - 1) / int64(sizePage), // total all page
		},
	}
	response.Response(ctx)
}

func (p *PayrollController) Payroll(ctx *gin.Context) {
	err := helper.Premission(ctx)
	if err != nil {
		return
	}

	var employee models.Employee
	var attedance []models.Attedance
	var req dto.RequestPayroll

	if err := dto.ValidationPayload(&req, ctx); err != nil {
		return
	}

	if err := p.DB.Model(&employee).Preload("Wallet").Where("id = ?", req.EmployeeId).First(&employee).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	datePattern := "%" + strings.ToLower(req.Date) + "%"
	if err := p.DB.Model(&attedance).Where("lower(chekin) LIKE ? AND employee_id = ?", datePattern, req.EmployeeId).Find(&attedance).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	salary_monthly := employee.Wallet.MonthlySalary
	attedances := attedance

	result_payroll := req.CalculationPayroll(salary_monthly, attedances)

	if err := p.DB.Create(&result_payroll).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    201,
		Message: "Payroll berhasil disimpan",
		Data:    result_payroll,
	}
	response.Response(ctx)
}

func (p *PayrollController) EmailPayslip(ctx *gin.Context) {
	err := helper.Premission(ctx)
	if err != nil {
		return
	}

	file, err := ctx.FormFile("pdf")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve PDF file",
		})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open PDF file",
		})
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read PDF file",
		})
		return
	}

	date := ctx.PostForm("date")

	go helper.SendPayslip(fileBytes, "arsatteguh@gmail.com", date)

	response := &helper.WithoutData{
		Code:    201,
		Message: "Payroll berhasil dikirim",
	}
	response.Response(ctx)
}

type payloadExcel struct {
	Id           int64     `json:"id"`
	EmployeeName string    `json:"employee_name"`
	DailySalary  float64   `json:"daily_salary"`
	Absence      int       `json:"absence"`
	Bonus        float64   `json:"bonus"`
	Status       string    `json:"status"`
	Tax          float64   `json:"tax"`
	TotalHour    float64   `json:"total_hour"`
	Total        float64   `json:"total"`
	Created      time.Time `json:"created" time_format:"2006-01-02"`
}

func (p *PayrollController) ExportExcelHandler(ctx *gin.Context) {
	var body []payloadExcel

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "payroll"
	index, _ := f.NewSheet(sheetName)

	f.SetActiveSheet(index)
	// Styling for header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "e4e4e7", // White text
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"0b192c"}, // Dark blue background
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Data style with borders
	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set column widths
	colWidths := []float64{5, 10, 25, 25, 10, 20, 10, 25, 15, 25, 20}
	for i, width := range colWidths {
		colName := string(rune('A' + i))
		f.SetColWidth(sheetName, colName, colName, width)
	}

	// Headers with improved readability
	headers := []string{
		"No",
		"ID",
		"Employee Name",
		"Daily Salary",
		"Absence",
		"Bonus",
		"Tax",
		"Total Hours",
		"Status",
		"Total",
		"Created Date",
	}

	// Set headers with styling
	for i, header := range headers {
		colName := string(rune('A' + i))
		f.SetCellValue(sheetName, colName+"1", header)
		f.SetCellStyle(sheetName, colName+"1", colName+"1", headerStyle)
	}

	// Write data rows with styling
	for i, payroll := range body {
		rowNum := i + 2
		rowData := []interface{}{
			i + 1,
			payroll.Id,
			payroll.EmployeeName,
			payroll.DailySalary,
			payroll.Absence,
			payroll.Bonus,
			payroll.Tax,
			payroll.TotalHour,
			payroll.Status,
			payroll.Total,
			payroll.Created.Format("2006-01-02"), // Format date consistently
		}

		for j, value := range rowData {
			colName := string(rune('A' + j))
			cellRef := fmt.Sprintf("%s%d", colName, rowNum)
			f.SetCellValue(sheetName, cellRef, value)
			f.SetCellStyle(sheetName, cellRef, cellRef, dataStyle)
		}
	}

	// Set headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=payroll.xlsx")

	// Write to response
	if err := f.Write(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
