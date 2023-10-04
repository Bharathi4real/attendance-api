document.getElementById('apiButton').addEventListener('click', () => {
    const url = "http://34.125.177.211:8080/login"; // Replace with your API endpoint URL
    const requestData = {
        username: 'exampleUsername',
        password: 'examplePassword'
    };
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('API Response:', data);
            // Handle the API response as needed
        })
        .catch(error => {
            console.error('API Request Error:', error);
            // Handle the error here
        });
});
