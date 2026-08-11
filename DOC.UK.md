# openrouter - довідник

Повний довідник пакета `openrouter`: клієнт, спільна модель `goloop/ai`, chat
completions (інтерфейс і нативний), стрімінг, атрибуція застосунку і моделі.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Створення клієнта](#створення-клієнта)
- [Generate і Stream](#generate-і-stream)
- [Структурований вивід](#структурований-вивід)
- [Пошук на боці провайдера](#пошук-на-боці-провайдера)
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

## Структурований вивід

`ai.Request.Format` лягає на власний `response_format` провайдера, тож запит на JSON
провайдер **дотримує**, а не просто «чує»:

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    "the-model",
	Messages: []ai.Message{ai.UserText("Склади SEO-поля для цієї статті.")},
	Format: &ai.Format{
		Type:   ai.FormatJSONSchema,
		Name:   "seo",
		Schema: schema,
	},
})

var seo SEO
err = resp.JSON(&seo)
```

`ai.FormatJSON` іде як `{"type":"json_object"}`, а `ai.FormatJSONSchema` - як
`{"type":"json_schema", ...}`. У простому JSON-режимі до system-промпта ще
додається `ai.Format.Instruction()`: цей wire-формат відбиває `json_object`,
якщо в повідомленнях ніде немає слова «json». Ваш власний system-промпт
лишається, інструкція йде після нього; схемний режим промпт не чіпає.


`ai.Response.Format` дорівнює `ai.FormatNative`: цей провайдер дотримує кожну
форму, яку приймає. Які моделі підтримують схемний режим - справа провайдера;
таблиці можливостей тут немає, тож про непідтримувану пару скаже він сам.

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

## Пошук на боці провайдера

`ai.Request.Hosted` отримує `ai.ErrNoHosted` ще до відправлення запиту.

Цей провайдер маршрутизує до моделей, а не обслуговує одну, тож можливості
запиту залежать від маршруту за ним, а пошук тут - радше властивість
маршрутизації, ніж інструмент, який пропонують моделі. Відповіді за провайдера
загалом цей пакет дати не може, а відповідь на кожен маршрут була б таблицею
можливостей, хибною того тижня, коли маршрут зміниться.

Ця відмова - задокументована поведінка, а не мовчазна прогалина: відповідь, яку
дали без замовленого пошуку, виглядає точно так само, як та, що з пошуком, тож
голосна помилка - єдиний спосіб їх розрізнити. Якщо відповідь усе одно потрібна,
повторіть запит без `Hosted`.

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
