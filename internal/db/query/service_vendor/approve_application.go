package vendor_query

import (
	"errors"
	"nearbyassist/internal/db"
)

func ApproveApplication(applicationId int) error {
	tx, err := db.Connection.Beginx()
	if err != nil {
		return err
	}

	updateStatus := `
        UPDATE 
            Application 
        SET
            status = 'approved'
        WHERE
            id = ?
    `

	if _, err := tx.Exec(updateStatus, applicationId); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.New("Failed to approve application and rollback transaction")
		}

		return err
	}

	promoteVendor := `
        INSERT INTO
            Vendor (vendorId, job)
        VALUES
            (
                (SELECT applicantId FROM Application WHERE id = ?),
                (SELECT job FROM Application WHERE id = ?)
            )
    `

	if _, err := tx.Exec(promoteVendor, applicationId, applicationId); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.New("Failed to promote applicant to vendor and rollback transaction")
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.New("Failed to commit transaction and rollback")
		}

		return err
	}

	return nil
}
