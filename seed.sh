#! /bin/bash

echo "cleaning up db tables..."
make migrate-reset

echo "migrating tables..."
make migrate-up

echo "seeding test admin account..."
na-cli create -u admin -p admin -m admin@email.com -r admin -e -k ~/repos/capstone/key.txt -c 'root:secret@/nearbyassist'
na-cli create -u admin2 -p admin2 -m admin2@email.com -r admin -e -k ~/repos/capstone/key.txt -c 'root:secret@/nearbyassist'
