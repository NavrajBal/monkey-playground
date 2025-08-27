export interface CodeSample {
  id: string;
  title: string;
  description: string;
  code: string;
}

export const codeSamples: CodeSample[] = [
  {
    id: "fibonacci",
    title: "Fibonacci Sequence",
    description: "Classic recursive fibonacci implementation",
    code: `let fibonacci = fn(x) {
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      return 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

puts("Fibonacci 5:", fibonacci(5));
puts("Fibonacci 8:", fibonacci(8));
puts("Fibonacci 10:", fibonacci(10));
fibonacci(12);`,
  },
  {
    id: "factorial",
    title: "Factorial Function",
    description: "Calculate factorial using recursion",
    code: `let factorial = fn(n) {
  if (n < 2) {
    1
  } else {
    n * factorial(n - 1);
  }
};

puts("5 factorial =", factorial(5));
puts("7 factorial =", factorial(7));
factorial(10);`,
  },
  {
    id: "arrays",
    title: "Array Operations",
    description: "Working with arrays and built-in functions",
    code: `let numbers = [1, 2, 3, 4, 5];
puts("Original array:", numbers);
puts("Length:", len(numbers));
puts("First element:", first(numbers));
puts("Last element:", last(numbers));
puts("Rest:", rest(numbers));

let extended = push(numbers, 6);
puts("After push:", extended);

let mixed = [1, "hello", true, fn(x) { x * 2 }];
mixed;`,
  },
  {
    id: "closures",
    title: "Closures Example",
    description: "Demonstrating closures and lexical scoping",
    code: `let makeMultiplier = fn(factor) {
  fn(x) {
    x * factor
  }
};

let double = makeMultiplier(2);
let triple = makeMultiplier(3);

puts("Double 5 =", double(5));
puts("Triple 4 =", triple(4));

let makeAdder = fn(base) {
  fn(x) {
    base + x
  }
};

let addTen = makeAdder(10);
puts("10 plus 7 =", addTen(7));
addTen(25);`,
  },
  {
    id: "higher-order",
    title: "Higher-Order Functions",
    description: "Functions that work with other functions",
    code: `let map = fn(arr, func) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), push(accumulated, func(first(arr))))
    }
  };
  iter(arr, [])
};

let double = fn(x) { x * 2 };
let square = fn(x) { x * x };

let numbers = [1, 2, 3, 4, 5];
puts("Original:", numbers);
puts("Doubled:", map(numbers, double));
map(numbers, square);`,
  },
  {
    id: "hash-maps",
    title: "Hash Maps",
    description: "Working with key-value data structures",
    code: `let person = {
  "name": "Alice",
  "age": 30,
  "city": "New York"
};

puts("Person:", person);
puts("Name:", person["name"]);
puts("Age:", person["age"]);

let mixed_keys = {
  "string": "value1",
  42: "numeric key",
  true: "boolean key"
};

mixed_keys;`,
  },
  {
    id: "conditionals",
    title: "Conditionals & Logic",
    description: "If/else statements and boolean logic",
    code: `let max = fn(a, b) {
  if (a > b) {
    a
  } else {
    b
  }
};

let grade = fn(score) {
  if (score > 89) {
    "A"
  } else {
    if (score > 79) {
      "B"
    } else {
      if (score > 69) {
        "C"
      } else {
        "F"
      }
    }
  }
};

puts("Max of 15 and 23:", max(15, 23));
puts("Grade for 85:", grade(85));
puts("Grade for 92:", grade(92));
grade(65);`,
  },
];
