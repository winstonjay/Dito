http://www.bbcbasic.co.uk/bbcwin/tutorial/

## Program
_==============================================================================_

START RETRO-BASIC!

REM My first retro basic program!

result = 50 * 2
PRINT "50 * 2 =", result
WAIT 100
PRINT TAB(4,2), "Hello"

REM Using TRUE / FALSE
INPUT "Enter password " Password$
IF Password$="Super" THEN
  Supervisor %= TRUE
ELSE
  Supervisor %= FALSE
ENDIF
IF Supervisor% THEN
  REM Show main configuration screen
  PRINT "Welcome, master"
ELSE
  REM Show ordinary user's screen
  PRINT "What do you want now?"
ENDIF
END

REM Geography quiz
PRINT "What is the capital of France:"
PRINT "a) Paris"
PRINT "b) London"
PRINT "c) Madrid"
INPUT "Enter a,b or c: " Reply$
CASE Reply$ OF
  WHEN "A", "a": PRINT "Correct"
  WHEN "B", "b": PRINT "Sorry, that's England"
  WHEN "C", "c": PRINT "Sorry, that's Spain"
  OTHERWISE: PRINT "Sorry, invalid response"
ENDCASE

END!

_==============================================================================_

## Expected Output
_==============================================================================_

50 * 2 = 100


    Hello



_==============================================================================