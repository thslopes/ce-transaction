package app

func (a *App) SetupRoutes() {
	a.FiberApp.Get("/transactions", a.GetTransactionsHandler)
}
