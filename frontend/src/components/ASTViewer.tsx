import React, { useCallback, useMemo } from "react";
import {
  ReactFlow,
  type Node,
  type Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  Position,
  Handle,
  BackgroundVariant,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { useTheme } from "../contexts/ThemeContext";
import "./ASTViewer.css";

interface ASTNode {
  type: string;
  [key: string]: any;
}

interface ASTViewerProps {
  astData: ASTNode | null;
}

// Custom node component for AST nodes
const ASTNodeComponent: React.FC<{ data: any }> = ({ data }) => {
  const getNodeColor = (type: string) => {
    switch (type) {
      case "Program":
        return "#3b82f6"; // Blue
      case "LetStatement":
        return "#10b981"; // Green
      case "ReturnStatement":
        return "#f59e0b"; // Yellow
      case "ExpressionStatement":
        return "#8b5cf6"; // Purple
      case "IfExpression":
        return "#ef4444"; // Red
      case "FunctionLiteral":
        return "#06b6d4"; // Cyan
      case "CallExpression":
        return "#f97316"; // Orange
      case "InfixExpression":
        return "#84cc16"; // Lime
      case "PrefixExpression":
        return "#ec4899"; // Pink
      case "Identifier":
        return "#6b7280"; // Gray
      case "IntegerLiteral":
        return "#14b8a6"; // Teal
      case "Boolean":
        return "#a855f7"; // Violet
      case "StringLiteral":
        return "#22c55e"; // Emerald
      default:
        return "#64748b"; // Slate
    }
  };

  return (
    <div
      className="ast-node"
      style={{ backgroundColor: getNodeColor(data.type) }}
    >
      <Handle type="target" position={Position.Top} />
      <div className="ast-node-content">
        <div className="ast-node-type">{data.type}</div>
        {data.value && <div className="ast-node-value">{data.value}</div>}
      </div>
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
};

const nodeTypes = {
  astNode: ASTNodeComponent,
};

const ASTViewer: React.FC<ASTViewerProps> = ({ astData }) => {
  const { theme } = useTheme();

  const { nodes, edges } = useMemo(() => {
    if (!astData) return { nodes: [], edges: [] };

    const nodes: Node[] = [];
    const edges: Edge[] = [];
    let nodeId = 0;

    const createNode = (node: any, parentId?: string, x = 0, y = 0): string => {
      const currentId = `node-${nodeId++}`;

      // Extract display value
      let value = "";
      if (node.Value !== undefined) value = String(node.Value);
      else if (node.Token?.Literal) value = node.Token.Literal;
      else if (node.Operator) value = node.Operator;
      else if (node.string) value = node.string;

      nodes.push({
        id: currentId,
        type: "astNode",
        position: { x, y },
        data: {
          type: node.type || "Unknown",
          value: value,
        },
      });

      if (parentId) {
        edges.push({
          id: `edge-${parentId}-${currentId}`,
          source: parentId,
          target: currentId,
        });
      }

      // Recursively create child nodes
      let childX = x - 100;
      let childY = y + 100;

      // Handle different AST node structures from Go backend
      if (node.statements && Array.isArray(node.statements)) {
        // Program node
        node.statements.forEach((stmt: any, index: number) => {
          createNode(stmt, currentId, childX + index * 250, childY);
        });
      }

      if (node.Name) {
        // LetStatement has Name and Value
        createNode(node.Name, currentId, childX, childY);
        childX += 200;
      }

      if (node.Value && node.Name) {
        // LetStatement Value
        createNode(node.Value, currentId, childX, childY);
      } else if (node.Value && !node.Name) {
        // Other Value nodes
        createNode(node.Value, currentId, x, childY);
      }

      if (node.Left) {
        createNode(node.Left, currentId, childX, childY);
        childX += 200;
      }

      if (node.Right) {
        createNode(node.Right, currentId, childX, childY);
      }

      if (node.Expression) {
        createNode(node.Expression, currentId, x, childY);
      }

      if (node.Condition) {
        createNode(node.Condition, currentId, childX, childY);
        childX += 200;
      }

      if (node.Consequence) {
        createNode(node.Consequence, currentId, childX, childY);
        childX += 200;
      }

      if (node.Alternative) {
        createNode(node.Alternative, currentId, childX, childY);
      }

      if (node.Parameters && Array.isArray(node.Parameters)) {
        node.Parameters.forEach((param: any, index: number) => {
          createNode(param, currentId, childX + index * 150, childY);
        });
      }

      if (node.Body) {
        createNode(node.Body, currentId, x, childY);
      }

      return currentId;
    };

    createNode(astData, undefined, 400, 50);
    return { nodes, edges };
  }, [astData]);

  const [flowNodes, setNodes, onNodesChange] = useNodesState(nodes);
  const [flowEdges, setEdges, onEdgesChange] = useEdgesState(edges);

  // Update nodes when astData changes
  React.useEffect(() => {
    setNodes(nodes);
    setEdges(edges);
  }, [nodes, edges, setNodes, setEdges]);

  const onConnect = useCallback(() => {}, []);

  if (!astData) {
    return (
      <div className="ast-viewer-empty">
        <p>No AST data to display. Parse some code first!</p>
      </div>
    );
  }

  return (
    <div className="ast-viewer-container">
      <ReactFlow
        nodes={flowNodes}
        edges={flowEdges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        fitView
        className={`ast-flow ${theme}`}
      >
        <Controls />
        <Background variant={BackgroundVariant.Dots} gap={20} size={1} />
      </ReactFlow>
    </div>
  );
};

export default ASTViewer;
