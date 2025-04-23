start_test_server:
	node index.js & \
	timeout 15 sh -c 'until nc -z localhost 3000; do sleep 5; done' || (echo "Test server failed to start." && exit 1);
