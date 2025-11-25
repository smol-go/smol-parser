# smol-parser

A JSON parser built from scratch in Go.

## Project Structure

```
smol-parser/
├── main.go         # Core parser implementation
├── main_test.go    # Comprehensive test suite
├── go.mod          # Go module file
└── README.md       # This file
```

## How It Works

### 1. Lexical Analysis (Lexer)

The lexer breaks the input string into tokens. It's like reading words in a sentence.

**Key Concepts:**
- **Token**: A meaningful unit (e.g., `{`, `"hello"`, `123`, `true`)
- **Scanning**: Reading one character at a time
- **Lookahead**: Peeking at the next character without consuming it

```go
type Lexer struct {
    input string  // The JSON string to parse
    pos   int     // Current position in input
    ch    byte    // Current character
}
```

**How it works:**
1. `readChar()` advances to the next character
2. `skipWhitespace()` ignores spaces, tabs, newlines
3. `NextToken()` identifies and returns the next token

**Example:**
```
Input: {"name": "John"}
Tokens: { → STRING("name") → : → STRING("John") → }
```

### 2. Syntactic Analysis (Parser)

The parser takes tokens and builds data structures. It uses **recursive descent parsing**.

**Key Concepts:**
- **Recursive descent**: Each grammar rule becomes a function
- **Current token**: The token we're looking at
- **Advance**: Move to the next token

```go
type Parser struct {
    lexer    *Lexer
    curToken Token  // Current token being examined
}
```

**Grammar Rules (simplified):**
```
Value  → Object | Array | String | Number | Boolean | Null
Object → { } | { Members }
Array  → [ ] | [ Elements ]
```

**How it works:**
1. `parseValue()` decides what type of value to parse
2. `parseObject()` handles `{...}` structures
3. `parseArray()` handles `[...]` structures

### 3. Parsing Flow Example

Let's trace: `{"name": "John", "age": 30}`

```
1. parseValue() sees { → calls parseObject()
2. parseObject():
   - Advance past {
   - See "name" (string key)
   - Advance, expect :
   - Call parseValue() → returns "John"
   - Store in map: {"name": "John"}
   - See , → continue
   - See "age" (string key)
   - Advance, expect :
   - Call parseValue() → returns 30.0
   - Store in map: {"name": "John", "age": 30.0}
   - See } → return map
```

## Key Implementation Details

### String Parsing

Handles escape sequences:
- `\"` → quote
- `\\` → backslash
- `\n` → newline
- `\uXXXX` → Unicode character

```go
// Input: "Hello\nWorld"
// Output: Hello
//         World
```

### Number Parsing

Supports full JSON number spec:
- Integers: `123`, `-456`
- Decimals: `123.456`
- Scientific: `1.5e-10`, `1E+10`

### Error Handling

Errors are returned with context:
```go
return nil, fmt.Errorf("expected colon after key")
```

## Usage

### Initialize Module

```bash
mkdir smol-parser
cd smol-parser
go mod init github.com/smol-go/smol-parser
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    jsonStr := `{"name": "Alice", "age": 30}`
    
    result, err := Parse(jsonStr)
    if err != nil {
        log.Fatal(err)
    }
    
    // Result is map[string]interface{}
    obj := result.(map[string]interface{})
    fmt.Println(obj["name"]) // Alice
    fmt.Println(obj["age"])  // 30
}
```

### Running Tests

```bash
go test -v
```

### Running Benchmarks

```bash
go test -bench=.
```

## Implementation Challenges & Solutions

### Challenge 1: String Escape Sequences

**Problem**: How to handle `\n`, `\t`, `\uXXXX`?

**Solution**: Switch statement in `readString()` with special handling for Unicode escapes. Read 4 hex digits, parse to int, convert to rune.

### Challenge 2: Number Parsing

**Problem**: JSON numbers can be complex: `-123.456e-10`

**Solution**: State machine approach:
1. Optional minus
2. Integer part (0 or 1-9 followed by digits)
3. Optional decimal point + digits
4. Optional exponent (e/E, optional +/-, digits)

### Challenge 3: Recursive Structures

**Problem**: Objects and arrays can contain each other infinitely

**Solution**: Recursive descent - `parseValue()` calls `parseObject()` which calls `parseValue()` again for nested values.

### Challenge 4: Error Recovery

**Problem**: How to report meaningful errors?

**Solution**: Track position in token, return descriptive error messages with context.

## Data Type Mapping

| JSON Type | Go Type              |
|-----------|----------------------|
| object    | map[string]interface{} |
| array     | []interface{}        |
| string    | string               |
| number    | float64              |
| boolean   | bool                 |
| null      | nil                  |

## Limitations

1. **Performance**: This parser prioritizes clarity over speed
2. **Numbers**: All numbers become float64 (JSON spec doesn't distinguish int/float)
3. **Big Numbers**: Very large integers may lose precision
4. **Memory**: Large JSON files load entirely into memory

## Learning Resources

**Concepts Demonstrated:**
- Lexical analysis
- Recursive descent parsing
- State machines
- Go interfaces (`interface{}` for dynamic types)
- Error handling patterns
- Table-driven tests

**Next Steps:**
- Add streaming parser (don't load entire file)
- Support custom struct unmarshaling
- Add JSON schema validation
- Implement JSON pointer (RFC 6901)
- Pretty printing/formatting

## Testing

The test suite covers:
- All JSON primitive types
- Nested structures
- Edge cases (empty arrays/objects)
- Error conditions
- Whitespace handling
- Escape sequences
- Benchmarks

Run specific tests:
```bash
go test -run TestParseString
go test -run TestParseObject
go test -bench=BenchmarkParseArray
```

## References

- [JSON Specification (RFC 8259)](https://tools.ietf.org/html/rfc8259)
- [Recursive Descent Parsing](https://en.wikipedia.org/wiki/Recursive_descent_parser)
- Go's official `encoding/json` package (for comparison)