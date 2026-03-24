package repo

import (
	"context"
	"database/sql"
	"producer-payment-notif/db"
	"producer-payment-notif/models"
	"time"
)

func GetDataQueueError() (res string, resError models.Respons, err error) {
	db, errcon := db.GetsSQLsrvDB()

	if errcon != nil {
		resError = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "Terkendala Jaringan/Koneksi, Silahkan Coba Beberapa Saat Lagi",
			Data:              nil,
		}

		return res, resError, errcon
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
		return res, resError, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "select json_req_error from trx_queue_error_notif where id=?", 3).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, resError, err
		}

		if ctx.Err() == context.DeadlineExceeded {
			resError = models.Respons{
				ResponseCode:      "408",
				ResponseMessage:   "System Timeout, Terkendala Jaringan /Koneksi, Silahkan Coba Beberapa Saat Lagi",
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            "System Timeout, Terkendala Jaringan /Koneksi, Silahkan Coba Beberapa Saat Lagi",
				Data:              nil,
			}
			return res, resError, err
		}

		resError = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Terjadi Kendala System (1001)",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "Terjadi Kendala System (1001)",
			Data:              nil,
		}
		return res, resError, err
	}

	tx.Commit()

	return res, resError, nil
}
