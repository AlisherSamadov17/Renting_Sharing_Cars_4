package postgres


import (
	"database/sql"
	"fmt"
	"rent-car/models"
	"rent-car/pkg"

	"github.com/google/uuid"
)

type carRepo struct {
	db *sql.DB
}

func NewCar(db *sql.DB) carRepo {
	return carRepo{
		db: db,
	}
}

func (c *carRepo) Create(car models.Car) (string, error) {

	id := uuid.New()

	query := ` INSERT INTO cars (
		id,
		name,
		brand,
		model,
		hourse_power,
		colour,
		engine_cap,
		year)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8) 
	`

	_,err := c.db.Exec(query,
		id.String(),
		car.Name, car.Brand,
		car.Model, car.HoursePower,
		car.Colour, car.EngineCap, car.Year)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (c *carRepo) Update(car models.Car) (string, error) {

	query := ` UPDATE cars set
			name=$1,
			brand=$2,
			model=$3,
			hourse_power=$4,
			colour=$5,
			engine_cap=$6,
			updated_at=CURRENT_TIMESTAMP
		WHERE id = $7 AND deleted_at = 0
	`

	_, err := c.db.Exec(query,
		car.Name, car.Brand,
		car.Model, car.HoursePower,
		car.Colour, car.EngineCap, car.Id)

	if err != nil {
		return "", err
	}

	return car.Id, nil
}

func (c *carRepo) GetAll(req models.GetAllCarsRequest) (models.GetAllCarsResponse, error) {
	var (
		resp   = models.GetAllCarsResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` and name ILIKE  '%%%v%%' `, req.Search)
	}

	filter += fmt.Sprintf("OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter:", filter)
	rows, err := c.db.Query(`select 
				count(id) OVER(),
				id, 
				name,
				brand,
				model,
				year,
				hourse_power,
				colour,
				engine_cap,
				created_at::date,
				updated_at,
				year
	  FROM cars WHERE deleted_at = 0 ` + filter + ``)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var (
			car      = models.Car{}
			updateAt sql.NullString
		)

		if err := rows.Scan(
			&resp.Count,
			&car.Id,
			&car.Name,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.HoursePower,
			&car.Colour,
			&car.EngineCap,
			&car.CreatedAt,
			&updateAt,
			&car.Year); err != nil {
			return resp, err
		}

		car.UpdatedAt = pkg.NullStringToString(updateAt)
		resp.Cars = append(resp.Cars, car)
	}
	return resp, nil
}

func (c *carRepo) GetByID(id string) (models.Car, error) {
	car := models.Car{}
	if err := c.db.QueryRow(`select id,name,brand,model,hourse_power,colour,engine_cap,created_at,year from cars where id = $1`, id).Scan(
		&car.Id,
		&car.Name,
		&car.Brand,
		&car.Model,
		&car.HoursePower,
		&car.Colour,
		&car.EngineCap,
		&car.CreatedAt,
		&car.Year,
	); err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (c *carRepo) Delete(id string) error {

	query := ` delete from cars WHERE id = $1`
	_, err := c.db.Exec(query, id)
	
	if err != nil {
		return err
	}

	return nil
}
