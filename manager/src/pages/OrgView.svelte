<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';

    export var params;

    function fetchOrg(id) {
        if (!params.id) {
            throw 'No org specified';
        }
        return gql({query: `{org (id: "${id}") {id, name, description, created}}`});
    }

    function onEdit() {
        push(`/org-edit/${params.id}`)
    }

</script>

<h1 class="title">View Org</h1>

{#await fetchOrg(params.id)}

    <progress class="progress is-small is-primary" max="100">15%</progress>
{:then d}
    <div class="content">
        <ul>
            <li>Name: {d.org.name}</li>
            <li>Description: {d.org.description}</li>
            <li>Created: {d.org.created}</li>
        </ul>
    </div>
    <button class="button" on:click={onEdit}>Edit</button>
{:catch error}
    <ErrorBar error={error}/>
{/await}
