start_test_server:
	node index.js & \
	timeout 5 sh -c 'until nc -z localhost 3000; do sleep 2; done' || (echo "Test server failed to start." && exit 1);
