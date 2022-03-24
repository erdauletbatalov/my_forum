package web

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"runtime/debug"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Извлекаем соответствующий набор шаблонов из кэша в зависимости от названия страницы
	// (например, 'home.page.html'). Если в кэше нет записи запрашиваемого шаблона, то
	// вызывается вспомогательный метод serverError(), который мы создали ранее.
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("шаблона %s не существует", name))
		return
	}

	// Рендерим файлы шаблона, передавая динамические данные из переменной `td`.
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validTags(tags []string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	for _, val := range tags {
		if re.FindStringSubmatch(val) == nil || !validInputStr(val) {
			return false
		}
	}
	return true
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func validInputStr(str string) bool {
	runes := []rune(str)
	if len(runes) < 3 || len(runes) > 30 {
		return false
	}
	return true
}

func validContent(str string) bool {
	runes := []rune(str)
	if len(runes) < 3 || len(runes) > 10000 {
		return false
	}
	return true
}
