<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';

    export var params;

    let name = '';
    let description = '';
    let error = null;
    let fetching = false;
    let success = null;

    onMount(async () => {
        if (params.id) {
            await fetchCustomer(params.id);
        } else {
            error = 'No customer specified';
        }
    })

    async function fetchCustomer(id) {
        fetching = true;
        error = false;
        name = '';
        description = '';

        try {
            let data = await gql({query: `{customer (name: "${id}") {name, description, created}}`});
            name = data.customer.name;
            description = data.customer.description;
        } catch(error) {
            error = 'Request failed (' + error + ')';
        }

        fetching = false;
    }

    async function handleSubmit()
    {
        if (name.length === 0) {
            error = 'No name specified'
            return
        }

        fetching = true;
        error = false;
        success = false;

        try {
            let data = await gql({query: `mutation {updateCustomer(name: "${params.id}", newName: "${name}") {name}}`});
            success = 'Customer successfully updated'
        } catch(e) {
            if (e instanceof Array) {
                e = e.map((err) => err.message).join(', ');
            }
            error = 'Request failed (' + e + ')';
        }

        fetching = false;
    }
</script>

<style>
form { width: 24rem;}
</style>

<h1 class="title">Edit Customer {params.id}</h1>

{#if error}<div class="notification is-danger has-text-centered">{error}</div>{/if}
{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

<form on:submit|preventDefault={handleSubmit}>

    <div class="field">
        <p class="control">
            <input bind:value={name} class="input {name.length === 0 && "is-danger"}" placeholder="Customer name">
        </p>
    </div>

    <div class="field">
        <p class="control">
            <textarea bind:value={description} class="textarea" placeholder="Customer description"/>
        </p>
    </div>

    <div class="field">
        <p class="control">
            <button class="button is-success">Save</button>
        </p>
    </div>

</form>

