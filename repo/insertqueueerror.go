package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"producer-payment-notif/db"
	"producer-payment-notif/models"
	"time"
)

func InsertQueueError(res []interface{}) (resError models.Respons, err error) {
	// resError := models.Respons{}

	db, errcon := db.GetsSQLsrvDB()

	if errcon != nil {
		resError = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			Data:              nil,
		}
		return resError, errcon
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resError = models.Respons{
			ResponseCode:      "408",
			ResponseMessage:   "System Timeout, Terkendala Jaringan /Koneksi, Silahkan Coba Beberapa Saat Lag",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "System Timeout, Terkendala Jaringan /Koneksi, Silahkan Coba Beberapa Saat Lagi",
			Data:              nil,
		}
		return resError, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()

	jsonData, err := json.Marshal(res)
	if err != nil {
		return resError, err
	}

	datas := string(jsonData)

	err = tx.QueryRowContext(ctx, "exec [spa_insert_queue_notif] @pJsonReqError=?, @pTemplateCodeNew=?, @pTemplateCodeOld=?", datas, "", "").Scan(&resError.ResponseCode, &resError.ResponseMessage)
	if err != nil {
		if err == sql.ErrNoRows {
			return resError, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			resError = models.Respons{
				ResponseCode:      "408",
				ResponseMessage:   "System Timeout, Terkendala Jaringan /Koneksi, Silahkan Coba Beberapa Saat Lagi",
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            "System Timeout, Terkendala Jaringan /Koneksi, Silahkan Coba Beberapa Saat Lagi",
				Data:              nil,
			}
			return resError, err
		}

		resError = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Terjadi Kendala System (1001)",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "Terjadi Kendala System (1001)",
			Data:              nil,
		}
		return resError, err
	}

	tx.Commit()

	return resError, nil
}
