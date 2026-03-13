# Challenge 2: Card Identification

This challenge involves reading a sequence of card numbers from user input and identifying their types based on defined prefix and length rules.

## Code Explanation

1. **`CardType` Struct**: 
   A structure that holds the criteria for a card format. It has a `Prefix` slide containing valid starting numbers (e.g., "62") and a `Length` slice for valid digit counts (e.g., 16, 17, 18).

2. **`cardTypes` Map**:
   A global map storing known card associations string -> `CardType`. Currently, it registers validation criteria for "China UnionPay" and "Switch".

3. **Input Handling**:
   The script uses `bufio.NewScanner` over `os.Stdin` to safely capture multiple card inputs. Empty lines or `EOF` break the scan loop.
   Values from a full line (multiple numbers could be placed here) are parsed by splitting whitespaces `strings.Fields(line)` and appended into a string slide `input`.

4. **Validation Logic**:
   - The outer loop iterates over all inputted card numbers.
   - For each number, a flag `found` tracks matched records.
   - The process searches through `cardTypes`, matching if `strings.HasPrefix(card, prefix)`.
   - If a prefix matches, it cross-verifies against all lengths within `cardTypeValue.Length`.
   - Once successfully identified, it prints out the card type and terminates the inner searches (`break` on matches).
   - If `found == false` after evaluating all prefixes of known cards, it declares the card unrecognizable.
