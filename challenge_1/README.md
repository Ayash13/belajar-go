# Challenge 1: Currency Conversion

This challenge demonstrates how to use structs, maps, slices, and loops in Go to build a simple currency conversion tool.

## Code Explanation

1. **`ConversionResult` Struct**: 
   Used to structure the result of a conversion, holding the `Currency` name (string) and the converted `Amount` (float64).

2. **`rates` Map**:
   A map that acts as a lookup table for currency exchange rates relative to IDR (e.g., USD is 15000, EUR is 16000).

3. **Input Handling**:
   The program uses `fmt.Scanf("%f", &nominal)` to capture the user's input (nominal value in IDR) from the terminal.

4. **Conversion Logic**:
   - An empty slice `results` of type `[]ConversionResult` is created.
   - A `for ... range` loop iterates over the `rates` map.
   - For each currency, it calculates the conversion (`nominal / rate`), creates a `ConversionResult` struct, and appends it to the `results` slice.

5. **Displaying Results**:
   A final `for ... range` loop iterates over the `results` slice and prints the currency and amount formatted to 2 decimal places (`%.2f`).
