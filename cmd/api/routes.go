package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes returns the HTTP router configured with all handlers, route-specific middleware,
// and global middleware.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Defined handlers for 404 and 205 status code
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck route
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Metrics debugging route
	router.Handler(http.MethodGet, "/v1/observability/metrics", expvar.Handler())

	// DATABASE SCHEMA ROUTES

	// hotel routes
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id", app.requirePermission("hotel:read", app.showHotelHandler))
	router.HandlerFunc(http.MethodGet, "/v1/hotels", app.requirePermission("hotel:read", app.listHotelsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/hotels", app.requirePermission("hotel:write", app.createHotelHandler))
	router.HandlerFunc(http.MethodPut, "/v1/hotels/:id", app.requirePermission("hotel:write", app.updateHotelHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/hotels/:id", app.requirePermission("hotel:write", app.deleteHotelHandler))

	// department routes
	router.HandlerFunc(http.MethodGet, "/v1/departments/:dept_name", app.requirePermission("department:read", app.showDepartmentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/departments", app.requirePermission("department:read", app.listDepartmentsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/departments", app.requirePermission("department:write", app.createDepartmentHandler))
	router.HandlerFunc(http.MethodPut, "/v1/departments/:dept_name", app.requirePermission("department:write", app.updateDepartmentHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/departments/:dept_name", app.requirePermission("department:write", app.deleteDepartmentHandler))

	// room_type routes
	router.HandlerFunc(http.MethodGet, "/v1/room_types/:id", app.requirePermission("room_type:read", app.showRoomTypeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/room_types", app.requirePermission("room_type:read", app.listRoomTypesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/room_types", app.requirePermission("room_type:write", app.createRoomTypeHandler))
	router.HandlerFunc(http.MethodPut, "/v1/room_types/:id", app.requirePermission("room_type:write", app.updateRoomTypeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/room_types/:id", app.requirePermission("room_type:write", app.deleteRoomTypeHandler))

	// room routes
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number", app.requirePermission("room:read", app.showRoomHandler))
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms", app.requirePermission("room:read", app.listRoomsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms", app.requirePermission("room:write", app.createRoomHandler))
	router.HandlerFunc(http.MethodPut, "/v1/hotels/:id/rooms/:number", app.requirePermission("room:write", app.updateRoomHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/hotels/:id/rooms/:number", app.requirePermission("room:write", app.deleteRoomHandler))

	// employee routes
	router.HandlerFunc(http.MethodGet, "/v1/employees/:email", app.requirePermission("employee:read", app.showEmployeeByEmailHandler))
	router.HandlerFunc(http.MethodGet, "/v1/employees", app.requirePermission("employee:read", app.listEmployeesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/employees", app.requirePermission("employee:write", app.createEmployeeHandler))
	router.HandlerFunc(http.MethodPut, "/v1/employees/:email", app.requirePermission("employee:write", app.updateEmployeeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/employees/:id", app.requirePermission("employee:write", app.deleteEmployeeHandler))
	// housekeeper
	router.HandlerFunc(http.MethodGet, "/v1/housekeepers/:id", app.requirePermission("employee:read", app.showEmployeeByIDHandler))
	// activation token
	router.HandlerFunc(http.MethodPut, "/v1/activated/employees", app.activateEmployeeHandler)

	// housekeeping_task routes
	router.HandlerFunc(http.MethodGet, "/v1/housekeeping_tasks/:taskID", app.requirePermission("housekeeping_task:read", app.showHousekeepingTaskHandler))
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number/housekeeping_tasks", app.requirePermission("housekeeping_task:read", app.listHousekeepingTasksHandler))
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms/:number/housekeeping_tasks", app.requirePermission("housekeeping_task:write", app.createHousekeepingTaskHandler))
	router.HandlerFunc(http.MethodPut, "/v1/housekeeping_tasks/:taskID", app.requirePermission("housekeeping_task:write", app.updateHousekeepingTaskHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/housekeeping_tasks/:taskID", app.requirePermission("housekeeping_task:write", app.deleteHousekeepingTaskHandler))

	// maintenance_report routes
	router.HandlerFunc(http.MethodGet, "/v1/maintenance_reports/:reportID", app.requirePermission("maintenance_report:read", app.showMaintenanceReportHandler))
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:id/rooms/:number/maintenance_reports", app.requirePermission("maintenance_report:read", app.listMaintenanceReportsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/hotels/:id/rooms/:number/maintenance_reports", app.requirePermission("maintenance_report:write", app.createMaintenanceReportHandler))
	router.HandlerFunc(http.MethodPut, "/v1/maintenance_reports/:reportID", app.requirePermission("maintenance_report:write", app.updateMaintenanceReportHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/maintenance_reports/:reportID", app.requirePermission("maintenance_report:write", app.deleteMaintenanceReportHandler))

	// guest routes
	router.HandlerFunc(http.MethodGet, "/v1/guests/:passport_number", app.requirePermission("guest:read", app.showGuestHandler))
	router.HandlerFunc(http.MethodGet, "/v1/guests", app.requirePermission("guest:read", app.listGuestsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/guests", app.requirePermission("guest:write", app.createGuestHandler))
	router.HandlerFunc(http.MethodPut, "/v1/guests/:passport_number", app.requirePermission("guest:write", app.updateGuestHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/guests/:passport_number", app.requirePermission("guest:write", app.deleteGuestHandler))

	// reservation routes
	router.HandlerFunc(http.MethodGet, "/v1/reservations/:reservationID", app.requirePermission("reservation:read", app.showReservationHandler))
	router.HandlerFunc(http.MethodGet, "/v1/reservations", app.requirePermission("reservation:read", app.listReservationsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/reservations", app.requirePermission("reservation:write", app.createReservationHandler))
	router.HandlerFunc(http.MethodPut, "/v1/reservations/:reservationID", app.requirePermission("reservation:write", app.updateReservationHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/reservations/:reservationID", app.requirePermission("reservation:write", app.deleteReservationHandler))

	// registration routes
	router.HandlerFunc(http.MethodGet, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.requirePermission("registration:read", app.showRegistrationHandler))
	router.HandlerFunc(http.MethodGet, "/v1/registrations/:reservationID", app.requirePermission("registration:read", app.listRegistrationsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/registrations", app.requirePermission("registration:write", app.createRegistrationHandler))
	router.HandlerFunc(http.MethodPut, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.requirePermission("registration:write", app.updateRegistrationHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/registrations/:reservationID/:hotelID/:roomNumber", app.requirePermission("registration:write", app.deleteRegistrationHandler))

	// token routes
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// global middleware
	return app.requestLogger( // first middleware
		app.metrics(
			app.recoverPanic(
				app.enableCORS(
					app.rateLimit(
						app.gzip(
							app.authenticate(router), // last middleware
						),
					),
				),
			),
		),
	)
}
