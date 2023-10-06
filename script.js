document.getElementById('apiButton').addEventListener('click', () => {
    const url = 'http://localhost:8080/login'; // Replace with your API endpoint URL
    const requestData = {
        username: 'admin',
        password: 'admin123'
    };

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
        .then(response => response.json())
        .then(data => {
            console.log('API Response:', data);
            // Handle the API response as needed
        })
        .catch(error => console.error('API Request Error:', error));
});
