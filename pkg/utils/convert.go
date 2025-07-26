package utils

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func Int32ToString(id int32) string {
	return strconv.FormatInt(int64(id), 10)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ConvertUUID(uuidStr string) (pgtype.UUID, error) {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}, nil
}
