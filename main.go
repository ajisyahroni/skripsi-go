package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err =
		gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/database_dev_go?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Gagal Conect Ke Database")
	}
	db.AutoMigrate(&student{})
}

type (
	student struct {
		gorm.Model
		Nama        string `json:"nama"`
		Alamat      string `json:"alamat"`
		NoHp        string `json:"no_hp"`
		Kelas       string `json:"kelas"`
		StatusAktif int    `json:"status_aktif"`
	}
	transformedStudent struct {
		ID          uint   `json:"id"`
		Nama        string `json:"nama"`
		Alamat      string `json:"alamat"`
		NoHp        string `json:"no_hp"`
		Kelas       string `json:"kelas"`
		StatusAktif bool   `json:"status_aktif"`
	}
)

func cretedStudent(c *gin.Context) {
	var std transformedStudent
	var model student
	c.Bind(&std)
	validasi := validatorCreated(std)
	model = transferVoToModel(std)
	if validasi != "" {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": validasi})
	} else {
		db.Create(&model)
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": model})
	}
}

func fetchAllStudent(c *gin.Context) {
	var model []student
	var vo []transformedStudent

	db.Find(&model)

	if len(model) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}

	for _, item := range model {
		vo = append(vo, transferModelToVo(item))
	}
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": vo})
}

func fetchSingleStuden(c *gin.Context) {
	var model student
	var vo transformedStudent

	modelID := c.Param("id")
	db.Find(&model, modelID)

	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}
	vo = transferModelToVo(model)
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": vo})
}

func updateStudent(c *gin.Context) {
	var model student
	var vo transformedStudent
	modelID := c.Param("id")
	db.First(&model, modelID)

	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}
	c.Bind(&vo)

	validasi := validatorCreated(vo)
	if validasi != "" {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": validasi})
	} else {
		db.Model(&model).Update(transferVoToModel(vo))
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": model})
	}
}

func deleteStudent(c *gin.Context) {
	var model student
	modelID := c.Param("id")

	db.First(&model, modelID)
	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Datat Tidak di Temukan"})
	}
	db.Delete(model)
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": "Data Telah berhasil di hapus"})
}

func transferModelToVo(model student) transformedStudent {
	var vo transformedStudent
	statusAktif := false
	if model.StatusAktif == 1 {
		statusAktif = true
	} else {
		statusAktif = false
	}
	vo = transformedStudent{
		ID:          model.ID,
		Nama:        model.Nama,
		Alamat:      model.Alamat,
		NoHp:        model.NoHp,
		StatusAktif: statusAktif,
		Kelas:       model.Kelas,
	}
	return vo
}

func transferVoToModel(vo transformedStudent) student {
	var model student
	statusAktif := 0
	if vo.StatusAktif == true {
		statusAktif = 1
	} else {
		statusAktif = 0
	}
	model = student{
		Nama:        vo.Nama,
		Kelas:       vo.Kelas,
		NoHp:        vo.NoHp,
		StatusAktif: statusAktif,
		Alamat:      vo.Alamat,
	}
	return model
}

func validatorCreated(vo transformedStudent) string {

	var kosong string = " Tidak Boleh Kosong"

	if vo.Nama == "" {
		return "Nama" + kosong
	}

	if vo.Alamat == "" {
		return "Alamat" + kosong
	}

	if vo.Kelas == "" {
		return "Kelas" + kosong
	}

	if vo.NoHp == "" {
		return "No Hp" + kosong
	}

	return ""
}

func main() {

	router := gin.Default()
	v1 := router.Group("/api/student")
	{
		v1.POST("", cretedStudent)
		v1.GET("", fetchAllStudent)
		v1.GET("/:id", fetchSingleStuden)
		v1.PUT("/:id", updateStudent)
		v1.DELETE("/:id", deleteStudent)
	}
	router.Run(":20001")
}
