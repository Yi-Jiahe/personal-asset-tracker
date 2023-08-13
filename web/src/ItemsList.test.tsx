import React from 'react';
import { findAllByText, render, screen } from '@testing-library/react';
import ItemsList from './ItemsList';
import axios from 'axios';

jest.mock("axios");
const axiosMocked = axios as jest.Mocked<typeof axios>;

test('Lists Items', () => {
  axiosMocked.get.mockResolvedValueOnce({
    "items": [
      {
        "item_id": 1234,
        "item_name": "potato"
      }
    ]
  });

  render(<ItemsList path="a/b/c" />);

  screen.findByText("potato");
});