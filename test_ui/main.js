// запустить npx http-server

const baseAddr = 'http://localhost:8091/api';

removeError = (errorDiv) => {
    errorDiv.innerHTML = '';
}

handleError = (errorStr) => {
    const errorDiv = document.getElementById('errors');
    errorDiv.innerHTML = `<p class="error">${errorStr}</p>`;
    setTimeout(removeError, 1000, errorDiv);
}


formData = () => {
    const fullInput = document.getElementById('full_input');
    const fullInputText =fullInput.value;
    if (fullInputText === '') {
        handleError('Текст не введен')
        return null;
    }

    const shortInput = document.getElementById('short_input');
    const shortInputText =shortInput.value;

    const linkObject = {
        real_url: fullInputText,
    };
    if (shortInputText !== '') {
        linkObject.shortcut = shortInputText;
    }
    return linkObject;
}

handleResult = (result) => {
    const hiddenDiv = document.getElementById("result");
    hiddenDiv.removeAttribute("hidden");
    const resultLink = document.getElementById("result_link");
    const linkText = baseAddr + '/link/' + result.shortcut;
    resultLink.href = linkText;
    resultLink.innerText = linkText;
    console.log(result);
}

formError = (status) => {
    console.log(status);
    switch (status) {
        case 403:
            return 'Это действие запрещено, короткая ссылка с таким именем уже создана';
        case 404:
            return 'Такой ссылки не существует';
        case 400:
            return 'Некорректный ввод';
        default:
            return `Ошибка с кодом ${status}`;

    }
}

function onButtonClick() {
    const data = formData();
    if (data == null) {
        console.log('couldn\'t form data');
        return;
    }
    fetch(baseAddr+'/link',
        {
                method: 'POST',
                mode: 'cors', // no-cors, *cors, same-origin
                cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
                credentials: 'include', // include, *same-origin, omit
                headers: {
                    'Content-Type': 'application/json',
                    'Charset': 'utf-8'

                },
                body: JSON.stringify(data)
    }).
        then(response => {return response;}).
        then(
            (data) => {
                if (data.status !== 201) {
                    handleError(formError(data.status));
                    console.log("failed to fetch");
                    return;
                }
                data.json().then(
                    handleResult
                )
                console.log(data);
            }
    )
}