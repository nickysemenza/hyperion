import React from 'react';
import renderer from 'react-test-renderer';
import { CueLabel } from './Cue';

test('Link changes the class when hovered', () => {
  const exampleCue = {
    status: 'active',
    elapsed_ms: 200,
    expected_duration: 400
  };
  const component = renderer.create(
    <CueLabel
      id={3}
      key={3}
      numActions={2}
      status={'ok'}
      cue={exampleCue}
      duration={600}
      duration_drift_ms={1}
      debug={false}
    />
  );
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();

  tree.props.debug = true;
  tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});
