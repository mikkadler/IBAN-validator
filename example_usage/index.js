document.getElementById('submitButton').addEventListener('click', function() {
    const inputText = document.getElementById('inputText').value;

    const inputArray = inputText.split('\n').filter(line => line.trim() !== '');;

    fetch('http://127.0.0.1:8080/validateIBAN', {
        method: 'POST',
        body: JSON.stringify({ data: inputArray }),
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(result => {
        const resultList = document.getElementById('result');
        resultList.innerHTML = ''; // Clear previous results
        result.forEach(item => {
            const listItem = document.createElement('li');
            listItem.innerText = "IBAN: " + item.iban + ", Valid: " + item.valid + ", Bank Name: " + item.bank_name + ", Error: " + item.error;
            resultList.appendChild(listItem);
        });
    })
    .catch(error => {
        console.error('Error:', error);
    });
});