#!/bin/bash
# start backend and save pid
APP_ENV=dev ./tupass-backend &
PID=$!
echo $PID
sleep 2s

# curl API and save response
RESPONSE=$(curl -H "password: \"test\"" -H "language: en" http://localhost:8000/api)
EXPECTED='{"length":{"score":15,"message":"very short","hint":"Your input has the length 4. That'"'"'s very bad."},"complexity":{"score":4,"message":"very simple","hint":"Your password has very few uppercase letters, digits and special characters."},"predictability":{"score":100,"message":"easy to predict","hint":"Your password is very similar to '"'test'"' in our password list."},"strength":{"score":9,"message":"very weak","hint":""}}'
echo "Got: $RESPONSE"

# test for expected result
(echo $RESPONSE | grep -q "$EXPECTED")

# save exit code / result and stop backend
RESULT=$?
kill $PID

# provide message and exit
if [ $RESULT -eq 0 ]
then
  echo "Response correct"
else
  echo "Expected: $EXPECTED"
  echo "Invalid respose"
fi

exit $RESULT
