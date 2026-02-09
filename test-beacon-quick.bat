@echo off
REM Quick test of beacon system

echo Testing Beacon System...
echo.

REM Test 1: Set beacon
echo Test 1: Setting beacon 'test-beacon'...
curl -s -X POST http://localhost:18080/api/beacon -H "Content-Type: application/json" -d "{\"name\":\"test-beacon\",\"metadata\":{\"test\":true}}"
echo.
echo.

REM Test 2: List beacons
echo Test 2: Listing beacons...
curl -s http://localhost:18080/api/beacons
echo.
echo.

REM Test 3: Read memory since beacon (should include the beacon itself)
echo Test 3: Reading memory since 'test-beacon'...
curl -s "http://localhost:18080/api/memory?beacon=test-beacon"
echo.
echo.

REM Test 4: Set another beacon
echo Test 4: Setting second beacon...
curl -s -X POST http://localhost:18080/api/beacon -H "Content-Type: application/json" -d "{\"name\":\"second-beacon\",\"metadata\":{\"note\":\"second test\"}}"
echo.
echo.

REM Test 5: Read memory from first to see both
echo Test 5: Reading memory since 'test-beacon' (should see both beacons)...
curl -s "http://localhost:18080/api/memory?beacon=test-beacon"
echo.
echo.

echo Test complete!
pause
