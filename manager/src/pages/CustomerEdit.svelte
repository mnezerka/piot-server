<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';

    export var params;

    let name = '';
    let description = '';
    let customer = null;
    let error = null;
    let fetching = false;
    let success = null;

    async function fetchCustomer() {
        if (!params.id) {
            error = 'No customer specified';
            return
        }

        fetching = true;
        error = null

        try {
            let data = await gql({query: `{customer (id: "${params.id}") {id, name, description, created}}`});
            customer = data.customer;
            name = customer.name;
            description = customer.description;
        } catch (e) {
            console.log('failed');
            error = e;
        }
        fetching = false;
    }

    async function updateCustomer()
    {
        if (name.length === 0) {
            error = 'No name specified'
            return
        }

        fetching = true;
        error = null;
        success = false;

        try {
            let data = await gql({query: `mutation {updateCustomer(id: "${params.id}", name: "${name}") {id}}`});
            success = 'Customer successfully updated'
        } catch(e) {
            error = e;
        }
        fetching = false;
    }

    onMount(fetchCustomer)

</script>

<style>
form { width: 24rem;}
</style>

<h1 class="title">Edit Customer</h1>

<ErrorBar error={error}/>

{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

{#if customer}
<form on:submit|preventDefault={updateCustomer}>

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
{/if}
