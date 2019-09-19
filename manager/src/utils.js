import {token} from './stores'

export async function gql(request) {

    let token_value = null;

    let unsubscribe = token.subscribe((value) => { token_value = value});
    unsubscribe();

    let response = await fetch(
        'http://localhost:9096/query',
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token_value,
            },
            body: JSON.stringify(request)
        })

    if (response.status !== 200) {
        throw response;
    }

    let data = await response.json();

    if (data.errors) {
        throw data.errors;
    }

    return data.data;
}

export function formatDate(timeStamp) {
    var d = new Date();
    d.setTime(timeStamp * 1000)
    return d.toUTCString();
}

