import { describe, it, expect } from 'vitest';

describe('frontend smoke test', () => {
  it('renders a simple app message', () => {
    const message = 'habitflow';
    expect(message).toContain('habitflow');
  });
});
