# Retry


There are 2 variable that trigger the retry mechanism
- Max Retry
  ErrorProducer < MaxRetry
  ErrorProducer == MaxRetry
  ErrorProducer > MaxRetry
  
- Error Type
    Any
    Specific and only one
    Specific and multiple in first
    Specific and multiple in middle
    Specific and multiple in last

Test:
- ErrorProducer < MaxRetry, ErrorType is Any
- ErrorProducer < MaxRetry, ErrorType is Specific and only one
- ErrorProducer < MaxRetry, ErrorType is Specific and multiple in first
- ErrorProducer < MaxRetry, ErrorType is Specific and multiple in middle
- ErrorProducer < MaxRetry, ErrorType is Specific and multiple in last
