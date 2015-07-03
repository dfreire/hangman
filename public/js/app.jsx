import React from 'react/addons';
import Router from 'react-router';
import {Route, Redirect, RouteHandler, Link} from 'react-router';

import Home  from './home.jsx'
import About from './about.jsx'

var pages = [
    { id: "home",  title: "Home",  handler: Home  },
    { id: "about", title: "About", handler: About }
];

pages.map(function(page, i) {
    page.path = "/" + page.id;
});

var activePath;

var App = React.createClass({
    render: function() {
        return (
            <div>
                <nav className="navbar navbar-default navbar-fixed-top">
                    <div className="container">
                        <div className="navbar-header">
                            <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
                                <span className="sr-only">Toggle navigation</span>
                                <span className="glyphicon glyphicon-menu-hamburger"></span>
                            </button>
                            <a className="navbar-brand" href="/">App</a>
                        </div>
                        <div id="navbar" className="collapse navbar-collapse">
                            <ul className="nav navbar-nav">
                                {pages.map(function(page, i) {
                                    var className = (activePath === page.path) ? "active" : "";
                                    return (
                                        <li key={i} className={className}><Link to={page.path}>{page.title}</Link></li>
                                    );
                                })}
                            </ul>
                        </div>
                    </div>
                </nav>
                <div className="container">
                    <div style={{marginTop: "70px"}} >
                        <RouteHandler/>
                    </div>
                </div>
            </div>
        )
    }
});

var routes = (
    <Route handler={App}>
        {pages.map(function(page, i) {
            return (<Route key={i} name={page.id} path={page.path} handler={page.handler}/>);
        })};
        <Redirect from="/" to="home" />
    </Route>
);

Router.run(routes, Router.HashLocation /*Router.HistoryLocation*/, function(Handler, state) {
    activePath = state.path.split("?")[0];
    React.render(<Handler />, document.getElementById('app'));
});
