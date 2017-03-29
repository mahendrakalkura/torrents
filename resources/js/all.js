var root = document.querySelector('#root');

var component = {
    view: function() {
        var node = m(
            'div',
            {
                class: 'progress',
            },
            m(
                'div', {
                    class: 'active progress-bar progress-bar-info progress-bar-striped',
                    style: 'widget: 0%',
                },
                '0%'
            )
        );
        return node;
    },
};

m.mount(root, component);
