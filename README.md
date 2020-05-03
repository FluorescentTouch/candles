# candles

## Makefile

This project is highly integrated with make tool, so to check all the possibilities:
```bash
make help
```

To run project:
```bash
make run
```

To run project with data race detector:
```bash
make run_race
```

To check project with static analyzers and tests:
```bash
make check
```

To check static checks only:
```bash
make static_check
```

To run tests (includes data race detector):
```bash
make test
```

## Helper tools

### For Go source code static analysis:

golangci-lint is used.

Install all tools:
```bash
make tools
```

Run all static checks:
```bash
make static_check
```
