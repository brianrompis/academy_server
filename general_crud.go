package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	function "github.com/brianrompis/academy_server/functions"
	"github.com/gorilla/mux"
)

func GetAllRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tableName := vars["table"]

	selectQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC", tableName)
	var responseData []map[string]interface{}
	rows, err := db.Raw(selectQuery).Rows()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		columns, err := rows.Columns()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			scanArgs := make([]interface{}, len(columns))
			values := make([]interface{}, len(columns))

			for i := range values {
				scanArgs[i] = &values[i]
			}

			err := rows.Scan(scanArgs...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			entry := make(map[string]interface{})
			for i, col := range values {
				columns[i] = function.SnakeToCamel(columns[i])
				if col != nil {
					entry[columns[i]] = col
				} else {
					entry[columns[i]] = nil
				}
			}

			responseData = append(responseData, entry)
		}
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func AddRecordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	tableName := vars["table"]

	columns := make([]string, 0)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		columns = append(columns, function.CamelToSnake(key))
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))

		formattedValue, ok := value.(string)
		if ok {
			values = append(values, formattedValue)
		} else {
			values = append(values, value)
		}
	}

	// insert created_at and updated_at timestamp
	now := time.Now()
	columns = append(columns, "created_at")
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
	values = append(values, now)
	columns = append(columns, "updated_at")
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
	values = append(values, now)

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	result := db.Exec(insertQuery, values...)
	if result.Error != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode("Data inserted successfully")
	}
}

func UpdateRecordHandler(w http.ResponseWriter, r *http.Request) {
	updateData := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	tableName := vars["table"]
	recordID, _ := strconv.Atoi(vars["id"])

	setValues := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range updateData {
		// modify updated_at timestamp
		snake_style_key := function.CamelToSnake(key)
		setValues = append(setValues, fmt.Sprintf("%s = $%d", function.CamelToSnake(key), len(values)+1))
		if snake_style_key != "updated_at" {
			values = append(values, value)
		} else {
			values = append(values, time.Now())
		}
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", tableName, strings.Join(setValues, ","), recordID)
	result := db.Exec(updateQuery, values...)
	if result.Error != nil {
		function.JsonResponse(w, http.StatusInternalServerError, result.Error.Error())
	} else {
		// fmt.Fprintf(w, fmt.Sprintf("%d row(s) affected.", result.RowsAffected))
		function.JsonResponse(w, http.StatusOK, fmt.Sprintf("Edit success. %d row(s) affected.", result.RowsAffected))
	}
}

func GetSingleRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tableName := vars["table"]
	recordID := vars["id"]

	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)
	var responseData []map[string]interface{}
	rows, err := db.Raw(selectQuery, recordID).Rows()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		columns, err := rows.Columns()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			scanArgs := make([]interface{}, len(columns))
			values := make([]interface{}, len(columns))

			for i := range values {
				scanArgs[i] = &values[i]
			}

			err := rows.Scan(scanArgs...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			entry := make(map[string]interface{})
			for i, col := range values {
				columns[i] = function.SnakeToCamel(columns[i])
				if col != nil {
					entry[columns[i]] = col
				} else {
					entry[columns[i]] = nil
				}
			}

			responseData = append(responseData, entry)
		}
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response headers and write the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tableName := vars["table"]
	recordID := vars["id"]

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	result := db.Exec(deleteQuery, recordID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "Record deleted successfully")
	}
}
