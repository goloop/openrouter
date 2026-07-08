# openrouter - довідник

Повний довідник пакета `openrouter`: клієнт, спільна модель `goloop/ai`, chat
completions (інтерфейс і нативний), стрімінг, атрибуція застосунку і моделі.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Створення клієнта](#створення-клієнта)
- [Generate і Stream](#generate-і-stream)
- [Нативні chat completions](#нативні-chat-completions)
- [Атрибуція застосунку](#атрибуція-застосунку)
- [Інструменти, зображення й system-промпти](#інструменти-зображення-й-system-промпти)
- [Моделі](#моделі)
- [Опції та помилки](#опції-та-помилки)

## Ментальна модель

`openrouter.Client` реалізує `ai.Client` - провайдер-незалежний контракт із
`github.com/goloop/ai`. OpenRouter - маршрутизувальний шлюз: один API і ключ
дають доступ до багатьох провайдерів моделей, де модель іменується як
`provider/model`. Спільні `Generate` і `Stream` покривають спільну основу (чат
із інструментами, зображеннями й стрімінгом), тож код проти інтерфейсу працює з
будь-яким провайдером.

Формат обміну - сумісний із chat completions; нативна частина - `ChatCompletion`
і перелік моделей.

```go
import (
	"github.com/goloop/ai"
	"github.com/goloop/openrouter"
)
```

## Створення клієнта

```go
c := openrouter.New(os.Getenv("OPENROUTER_API_KEY"))

c = openrouter.New(apiKey,
	openrouter.WithReferer("https://myapp.example"),
	openrouter.WithTitle("My App"),
)
```

Base URL за замовчуванням `https://openrouter.ai/api/v1`.

## Generate і Stream

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    openrouter.ModelClaudeSonnet,
	System:   "You are concise.",
	Messages: []ai.Message{ai.UserText("Name three primary colors.")},
})
resp.Text()
resp.ToolCalls()
resp.Usage
```

`Stream` повертає `iter.Seq2[ai.Chunk, error]`: текстові дельти чанками з `Text`,
завершений виклик інструмента - чанком із `ToolCall`, фінальний чанк - `Done` і
`Usage`.

## Нативні chat completions

Для опцій, специфічних для провайдера, будуйте `ChatRequest` і викликайте
`ChatCompletion` чи `ChatCompletionStream`:

```go
resp, err := c.ChatCompletion(ctx, &openrouter.ChatRequest{
	Model:          openrouter.ModelGPT4o,
	Messages:       []openrouter.ChatMessage{{Role: "user", Content: "as JSON"}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

## Атрибуція застосунку

`WithReferer` і `WithTitle` виставляють заголовки `HTTP-Referer` і `X-Title`,
якими OpenRouter атрибутує й ранжує трафік застосунку. Обидва необовʼязкові.

## Інструменти, зображення й system-промпти

Інструменти, зображення й system-промпти використовують спільні типи `ai`:
`ai.Tool`, `ai.Image`, `ai.ToolResult` і повідомлення `RoleSystem` або поле
`System`. Результати інструментів надсилаються назад повідомленнями `RoleTool`,
де `ai.ToolResult.ID` збігається з `ai.ToolUse.ID`. Вбудовані байти зображення
надсилаються як base64 data URI.

## Моделі

Працює будь-який рядок `provider/model`; для зручності є константи `ModelGPT4o`,
`ModelClaudeSonnet`, `ModelGeminiFlash` тощо.

```go
models, err := c.Models(ctx)
models[0].ID            // "openai/gpt-4o"
models[0].ContextLength
```

## Опції та помилки

Опції: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`,
`WithHeader`, `WithReferer`, `WithTitle`.

Невдала відповідь стає `*ai.APIError` зі `Status`, `Type`, `Code`, `Message` і
сирим тілом:

```go
var apiErr *ai.APIError
if errors.As(err, &apiErr) && apiErr.Status == http.StatusTooManyRequests {
	// backoff
}
```

Запити без моделі чи повідомлень падають до мережі з `ai.ErrNoModel` або
`ai.ErrNoMessages`.
