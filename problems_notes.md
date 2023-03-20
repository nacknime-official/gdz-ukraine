# Ideas and notes for solving problems

### 1. Checking the user's input (its validity and handling "…")
Because I'm going to cache the homework's API calls, I can just call the homework API for the previous step and check is the user's input is in the result. Smth like that:
```go
func (h *userHandler) OnInputYear(c telebot.Context, state fsm.Context) error {
    // get data from state
    ...

    years, err := h.homeworkService.GetYears(opts)
    if c.Message().Text not in years {
        return c.Send("The input is wrong")
        // and maybe even recreate the markup with the new data, because the old markup can be outdated
    }
}
```

Also we are solving the "…" problem that way (when the button's text is >128 symbols): we just check if the user's input ends with "…" symbol and search for the right entity in the result from the service:
```go
func (h *userHandler) OnInputAuthor(c telebot.Context, state fsm.Context) error {
    // get data from state
    ...

    authors, err := h.homeworkService.GetAuthors(opts)
    if c.Message().Text[-1] == "…" {
        if c.Message().Text not in authors that starts with c.Message().Text[:-1] {...}
    } else {
        if c.Message().Text not in authors {...}
    }
}
```
