package context

type ContextOptions struct {
    DbUri string
    DbName string
    LogLevel string
}

func NewContextOptions() *ContextOptions {
    o := &ContextOptions{
        DbUri:          "piot",
        DbName:         "piot",
        LogLevel:       "INFO",
    }
    return o
}

