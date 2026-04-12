require('babel-register')({
  presets: ['es2015']
});

const fs = require('fs');
const path = require('path');
const b3 = require('../../src');

function installFakeNow(now) {
  if (typeof now !== 'number') {
    return () => {};
  }

  const RealDate = Date;

  class FakeDate extends RealDate {
    constructor(...args) {
      if (args.length === 0) {
        super(now);
        return;
      }
      super(...args);
    }

    static now() {
      return now;
    }
  }

  global.Date = FakeDate;
  return () => {
    global.Date = RealDate;
  };
}

function applyBlackboardState(blackboard, treeId, state) {
  if (!state) {
    return;
  }

  if (state.base) {
    Object.keys(state.base).forEach((key) => {
      blackboard.set(key, state.base[key]);
    });
  }

  if (state.tree) {
    Object.keys(state.tree).forEach((key) => {
      blackboard.set(key, state.tree[key], treeId);
    });
  }

  if (state.nodes) {
    Object.keys(state.nodes).forEach((nodeId) => {
      const nodeState = state.nodes[nodeId];
      Object.keys(nodeState).forEach((key) => {
        blackboard.set(key, nodeState[key], treeId, nodeId);
      });
    });
  }
}

function normalizeDump(value) {
  return JSON.parse(JSON.stringify(value));
}

function runFixture(fixture) {
  const tree = new b3.BehaviorTree();
  tree.load(fixture.tree);

  const blackboard = new b3.Blackboard();
  applyBlackboardState(blackboard, tree.id, fixture.blackboard);

  const statuses = [];

  (fixture.ticks || []).forEach((tick) => {
    const restore = installFakeNow(tick.now);
    try {
      statuses.push(tree.tick(tick.target || null, blackboard));
    } finally {
      restore();
    }
  });

  const openNodes = (blackboard.get('openNodes', tree.id) || []).map((node) => node.id);
  const nodeCount = blackboard.get('nodeCount', tree.id);

  return {
    statuses,
    treeMemory: {
      openNodes,
      nodeCount
    },
    dump: normalizeDump(tree.dump())
  };
}

function runFixtures(fixturePaths) {
  const results = {};

  fixturePaths.forEach((fixturePath) => {
    const input = fs.readFileSync(path.resolve(fixturePath), 'utf8');
    const fixture = JSON.parse(input);
    results[fixture.name] = runFixture(fixture);
  });

  return results;
}

function main() {
  const fixturePaths = process.argv.slice(2);
  if (fixturePaths.length === 0) {
    throw new Error('at least one fixture path is required');
  }

  if (fixturePaths.length === 1) {
    const input = fs.readFileSync(path.resolve(fixturePaths[0]), 'utf8');
    const fixture = JSON.parse(input);
    process.stdout.write(JSON.stringify(runFixture(fixture)));
    return;
  }

  process.stdout.write(JSON.stringify(runFixtures(fixturePaths)));
}

main();
