var url = document.querySelector('body').getAttribute('data-url');

var root = document.querySelector('#root');

var oninit = function() {
    component.items = [];
    component.isLoading = true;
    component.isSuccess = false;
    component.isFailure = false;
    m
        .request({method: 'GET', url: url + '/items/'})
        .then(
            function(items) {
                component.items = items;
                component.isLoading = false;
                component.isSuccess = true;
                component.isFailure = false;
            },
            function() {
                component.items = [];
                component.isLoading = false;
                component.isSuccess = false;
                component.isFailure = true;
            }
        );
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
            return m(
                'table',
                {class: 'table table-bordered table-hover table-striped'},
                m('tbody', {}, component.items.map(function(item) {
                    return [
                        m(
                            'tr',
                            {},
                            [
                                m('td', {class: 'text-narrow'}, item.category),
                                m('td', {class: 'text-narrow text-right'}, item.seeds),
                                m('td', {}, m('a', {href: item.url}, item.title)),
                                m('td', {class: 'text-narrow text-center'}, m('a', {href: item.magnet}, 'Magnet')),
                            ]
                        ),
                        m(
                            'tr',
                            {},
                            m('td', {colspan: 4}, item.urls.map(function(url) {
                                return m('a', {class: 'block', href: url}, url);
                            }))
                        ),
                    ];
                })
            ));
        }
        if (component.isFailure) {
            return m('div', {class: 'alert alert-danger'}, 'An unknown error has occurred. Please try again.');
        }
    },
};

m.mount(root, component);
