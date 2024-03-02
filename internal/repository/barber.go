package repository

import (
	"adv_programming_3_4-main/internal/model"
	"database/sql"
	"strconv"
)

type BarberRepository struct {
	db *sql.DB
}

func NewBarberRepository(db *sql.DB) *BarberRepository {
	return &BarberRepository{db: db}
}

func (h *BarberRepository) GetBarbersFromDB() ([]model.Barber, error) {
	rows, err := h.db.Query("SELECT id, name, basic_info, price, experience, status, image_path FROM barbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []model.Barber
	for rows.Next() {
		var b model.Barber
		if err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath); err != nil {
			return nil, err
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return barbers, nil
}

func (h *BarberRepository) GetFilteredBarbersFromDB(statusFilter, experienceFilter, sortBy, pageStr string, itemsPerPage int) ([]model.Barber, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	query := `SELECT id, name, basic_info, price, experience, status, image_path FROM barbers WHERE true`
	if statusFilter != "" {
		query += " AND status = '" + statusFilter + "'"
	}
	if experienceFilter != "" {
		query += " AND experience = '" + experienceFilter + "'"
	}
	switch sortBy {
	case "name":
		query += " ORDER BY name"
	case "price":
		query += " ORDER BY price"
	}
	query += " LIMIT " + strconv.Itoa(itemsPerPage) + " OFFSET " + strconv.Itoa((page-1)*itemsPerPage)

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []model.Barber
	for rows.Next() {
		var b model.Barber
		if err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath); err != nil {
			return nil, err
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return barbers, nil
}
