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

    let data = await response.json();

    return data.data;
}
