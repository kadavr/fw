package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadavr/fw/crud"
)

type HotelWithImagesAndRooms struct {
	Hotel  `gorm:"embedded"`
	Images []Image `gorm:"foreignKey:HotelID;references:ID" json:"images"`
	Rooms  []Room  `gorm:"foreignKey:HotelID;references:ID" json:"rooms"`
}

type Hotel struct {
	crud.ModelBase
	Name    string    `json:"name,omitempty"`
	BeginAt time.Time `json:"begin_at,omitempty"`
	EndAt   time.Time `json:"end_at,omitempty"`
}

type Image struct {
	crud.ModelBase
	HotelID int64  `json:"hotel_id,omitempty"`
	URL     string `json:"url,omitempty"`
}

type Room struct {
	crud.ModelBase
	ID         int64  `json:"id,omitempty"`
	HotelID    int64  `json:"hotel_id,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`
	IsDeleted  bool   `json:"is_deleted,omitempty"`
	Name       string `json:"name,omitempty"`
	RoomTypeID int64  `json:"room_type_id,omitempty"`
}

type ModelDict struct {
	crud.ModelBase
	Alias string
	Name  string
}
type City struct {
	ModelDict //`gorm:"embedded;embeddedPrefix:dict_"`
}

func main() {
	router := gin.Default()

	//gormInstance := crud.LoadGormDB()

	//externRepo := crud.OptCrudRepo[Hotel](crud.ExtRepo[Hotel]{})
	_ = crud.Crud[Hotel](router, nil)
	_ = crud.Crud[Image](router, nil)
	_ = crud.Crud[Room](router, nil)

	crd := crud.Crud[HotelWithImagesAndRooms](router, nil)

	if err := crd.Migrate(&Hotel{}, &Image{}, &Room{}); err != nil {
		fmt.Println(err)
		return
	}
	if err := router.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}

// func main() {
// 	router := gin.Default()
// 	crd := crud.CrudNest[Hotel, Image](router, nil)
// 	crd.CRud.ReadManyBeforeFind = func(d *gorm.DB) *gorm.DB {
// 		return d.Where("1=1")
// 	}
// 	if err := crd.Migrate(); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("crd: %+v", crd)
// 	if err := router.Run(":8080"); err != nil {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	app := fiber.New(fiber.Config{
// 		JSONEncoder: json.Marshal,
// 		JSONDecoder: json.Unmarshal,
// 		//Prefork:     true,
// 	})
// 	crd := crud.CrudFiber[Hotel](app)
// 	if err := crd.Migrate(); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("crd: %+v", crd)
// 	if err := app.Listen(":8080"); err != nil {
// 		fmt.Println(err)
// 	}

// }
