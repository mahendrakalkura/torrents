jQuery(function() {
    var root = document.querySelector('#root');

    var component = m(
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

    m.render(root, component);
});
