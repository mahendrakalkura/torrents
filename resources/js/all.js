var url = document.querySelector('body').getAttribute('data-url');

var root = document.querySelector('#root');

var oninit = function() {
    component.categories = [];
    component.items = [];
    component.active = undefined;
    component.isLoading = true;
    component.isSuccess = false;
    component.isFailure = false;
    m
        .request({method: 'GET', url: url + '/items/'})
        .then(
            function(items) {
                component.items = items;
                component.categories = component.items.map(function(item) {
                    return item.category;
                });
                component.categories = new Set(component.categories);
                component.categories = Array.from(component.categories);
                component.categories.sort();
                component.active = component.categories[0];
                component.isLoading = false;
                component.isSuccess = true;
                component.isFailure = false;
            },
            function() {
                component.items = [];
                component.categories = [];
                component.active = undefined;
                component.isLoading = false;
                component.isSuccess = false;
                component.isFailure = true;
            }
        );
};

var categories = {
    view: function (component) {
        return m(
            'div',
            {class: 'btn-group'},
            component.attrs.categories.map(function(category) {
                var class_ = [];
                class_.push('btn');
                class_.push('btn-default');
                if (category === component.attrs.active) {
                    class_.push('active');
                }
                class_ = class_.join(' ');
                return m(
                    'a',
                    {
                        class: class_,
                        onclick: function() {
                            component.attrs.active = category;
                        }
                    },
                    category
                );
            })
        );
    }
};

var items = {
    view: function (component) {
        var items = component.attrs.items.filter(function (item) {
            return item.category === component.attrs.active;
        });
        return m(
            'div',
            {},
            m(
                'table',
                {class: 'table table-bordered table-hover table-striped'},
                m('tbody', {}, items.map(function(item) {
                    var trs = [];
                    trs.push(m(
                        'tr',
                        {},
                        [
                            m('td', {}, m('a', {href: item.url}, item.title)),
                            m('td', {class: 'text-narrow'}, item.category),
                            m('td', {class: 'text-narrow'}, item.timestamp),
                            m('td', {class: 'text-narrow text-right'}, item.seeds),
                            m('td', {class: 'text-narrow text-center'}, m('a', {href: item.magnet}, 'Magnet')),
                        ]
                    ));
                    if (item.urls.length > 0) {
                        trs.push(m(
                            'tr',
                            {},
                            m('td', {colspan: 5}, item.urls.map(function(url) {
                                return m('a', {class: 'block', href: url}, url);
                            }))
                        ));
                    }
                    return trs;
                })
            )
        ));
    }
};

var component = {
    items: [],
    isLoading: true,
    isSuccess: false,
    isFailure: false,
    oninit: oninit,
    view: function() {
        if (component.isLoading) {
            return m('p', {class: 'text-center'}, m('i', {class: 'fa fa-2x fa-cog fa-spin'}));
        }
        if (component.isSuccess) {
            return [m(categories, component), m(items, component)];
        }
        if (component.isFailure) {
            return m('div', {class: 'alert alert-danger'}, 'An unknown error has occurred. Please try again.');
        }
    },
};

m.mount(root, component);
