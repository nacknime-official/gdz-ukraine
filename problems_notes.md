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

### 2. Make the code in handlers clearer
I think about a kind of setting object or a runtime config that maps a state to all the info for "back" button (get the previous handler and message), checking input (callback for getting the data again and check the input with the data) etc.

Like something similar to aiogram-dialogs maybe.

### 3. Save IDs of the entities and pass the ID to the Homework gateway
The problem is that the bot doesn't use an inline keyboard, it uses a reply keyboard which doesn't have *callback data* with which I could link the item to its ID.

The bot can't use an inline keyboard because of the huge amount of the buttons and it looks bad with buttons' long texts.

Sure, I can pass only names to the gateway, like the bot's python version do.
But I found out another way to do that: after the user chosen the entity (a subject for example) I:
- get the subject list again
- find the chosen subject in the list
- get its ID
- and set the ID to the user's state instead of the value of the entity.

Anyway I get the same data in the next handler again to check that user "clicked on one of the buttons" and didn't input something else. Yes, I began to like this variant more after realising this fact.

But I should recheck if the Vshkole API has IDs for everything I need.

### 4. Amount of buttons can be more than 300
I've found out that the maximum amount of buttons for a reply keyboard is **300**.

And there are some steps that has >=300 choices.

So, it seems we need to implement pagination if the result from gateway is >=300.
