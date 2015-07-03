import React from "react/addons";
import Router from "react-router";

export default React.createClass({
    mixins: [Router.State],

    render: function() {
        return (<h1>Home</h1>);
    }
});
