<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';

    export var params;

    function fetchCustomer(id) {
        if (!params.id) {
            throw 'No customer specified';
        }
        return gql({query: `{customer (id: "${id}") {id, name, description, created}}`});
    }

    function onEdit() {
        push(`/customer-edit/${params.id}`)
    }

</script>

<h1 class="title">View Customer</h1>

{#await fetchCustomer(params.id)}

    <progress class="progress is-small is-primary" max="100">15%</progress>
{:then d}
    <div class="content">
        <ul>
            <li>Name: {d.customer.name}</li>
            <li>Description: {d.customer.description}</li>
            <li>Created: {d.customer.created}</li>
        </ul>
    </div>
    <button class="button" on:click={onEdit}>Edit</button>
{:catch error}
    <ErrorBar error={error}/>
{/await}
