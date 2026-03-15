import { clsx } from "clsx";
import { ButtonHTMLAttributes, forwardRef } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  size?: "sm" | "md" | "lg";
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = "primary", size = "md", className, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={clsx(
          "inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50",
          {
            "bg-indigo-600 text-white hover:bg-indigo-700": variant === "primary",
            "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50":
              variant === "secondary",
            "text-gray-600 hover:text-gray-900 hover:bg-gray-100":
              variant === "ghost",
            "bg-red-600 text-white hover:bg-red-700": variant === "danger",
            "px-2.5 py-1.5 text-sm": size === "sm",
            "px-4 py-2 text-sm": size === "md",
            "px-6 py-3 text-base": size === "lg",
          },
          className
        )}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";
